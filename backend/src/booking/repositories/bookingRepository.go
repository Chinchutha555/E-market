package repositories

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/booking/domains"
	"gorm.io/gorm"
)

type repository struct {
	conn *gorm.DB
}

func NewRepository(db *gorm.DB) domains.Repository {
	return &repository{conn: db}
}

// promoteOrReleaseStallOnTx checks if a stall has an active primary holder or active waitlist offer;
// if not, it re-indexes remaining waitlist records and offers the stall to position #1 (48h window)
// or sets the stall back to 'available' if no one is in the queue.
func promoteOrReleaseStallOnTx(tx *gorm.DB, stallIdent string) error {
	now := time.Now()
	cleanIdent := strings.TrimSpace(stallIdent)
	if cleanIdent == "" {
		return nil
	}

	// 1. Check if there is still an active primary booking holder on this stall
	var activeHolderCount int64
	if err := tx.Model(&models.Booking{}).
		Where("stall_ident = ? AND status = 'A' AND queue_number = 0 AND booking_status NOT IN ('waiting_queue', 'rejected', 'cancelled', 'not_interested', 'no_show', 'terminated')", cleanIdent).
		Count(&activeHolderCount).Error; err != nil {
		return err
	}
	if activeHolderCount > 0 {
		return nil
	}

	// 2. Re-index all active Waitlist items (offered + waiting) in order
	var activeWaitlist []models.Waitlist
	if err := tx.Where("stall_ident = ? AND status IN ('offered', 'waiting')", cleanIdent).
		Order("CASE WHEN status = 'offered' THEN 0 ELSE 1 END ASC, position ASC, id ASC").
		Find(&activeWaitlist).Error; err != nil {
		return err
	}

	for idx := range activeWaitlist {
		oldPos := activeWaitlist[idx].Position
		newPos := idx + 1
		if oldPos != newPos {
			if err := tx.Model(&activeWaitlist[idx]).Updates(map[string]interface{}{
				"position":   newPos,
				"updated_at": now,
			}).Error; err != nil {
				return err
			}
			activeWaitlist[idx].Position = newPos
		}
	}

	// 3. If the first item is already 'offered' and not expired, ensure stall is 'offered'
	if len(activeWaitlist) > 0 && activeWaitlist[0].Status == "offered" && activeWaitlist[0].ExpiresAt != nil && activeWaitlist[0].ExpiresAt.After(now) {
		waiterPos := 1
		for j := range activeWaitlist {
			if activeWaitlist[j].Status == "waiting" {
				if activeWaitlist[j].Position != waiterPos {
					if err := tx.Model(&activeWaitlist[j]).Updates(map[string]interface{}{
						"position":   waiterPos,
						"updated_at": now,
					}).Error; err != nil {
						return err
					}
					activeWaitlist[j].Position = waiterPos
				}
				waiterPos++
			}
		}
		return tx.Model(&models.Stall{}).Where("ident = ?", cleanIdent).Update("stall_status", "offered").Error
	}

	// 4. Otherwise, promote the first 'waiting' item to 'offered'
	for idx := range activeWaitlist {
		if activeWaitlist[idx].Status == "waiting" {
			nextWaiter := &activeWaitlist[idx]
			offerHours := 48
			var mkt models.Market
			if nextWaiter.MarketID != nil && *nextWaiter.MarketID > 0 {
				if tx.First(&mkt, *nextWaiter.MarketID).Error == nil && mkt.WaitlistOfferHours > 0 {
					offerHours = mkt.WaitlistOfferHours
				}
			} else if strings.TrimSpace(nextWaiter.MarketName) != "" {
				if tx.Where("name = ?", strings.TrimSpace(nextWaiter.MarketName)).First(&mkt).Error == nil && mkt.WaitlistOfferHours > 0 {
					offerHours = mkt.WaitlistOfferHours
				}
			}
			expires := now.Add(time.Duration(offerHours) * time.Hour)
			if err := tx.Model(nextWaiter).Updates(map[string]interface{}{
				"status":               "offered",
				"position":             1,
				"offered_at":           now,
				"expires_at":           expires,
				"near_expiry_notified": false,
				"updated_at":           now,
			}).Error; err != nil {
				return err
			}
			nextWaiter.Status = "offered"

			// Re-index remaining waiting items so the next in line moves up: position 1, 2, 3...
			waiterPos := 1
			for j := range activeWaitlist {
				if activeWaitlist[j].ID == nextWaiter.ID {
					continue
				}
				if activeWaitlist[j].Status == "waiting" {
					if activeWaitlist[j].Position != waiterPos {
						if err := tx.Model(&activeWaitlist[j]).Updates(map[string]interface{}{
							"position":   waiterPos,
							"updated_at": now,
						}).Error; err != nil {
							return err
						}
						activeWaitlist[j].Position = waiterPos
					}
					waiterPos++
				}
			}

			if err := tx.Model(&models.Stall{}).Where("ident = ?", cleanIdent).Update("stall_status", "offered").Error; err != nil {
				return err
			}

			mktLabel := strings.TrimSpace(nextWaiter.MarketName)
			if mktLabel == "" {
				mktLabel = "ตลาด กทม."
			}
			notifs := []models.Notification{
				{
					RecipientRole:   "tenant",
					RecipientUserID: &nextWaiter.TenantID,
					RecipientName:   nextWaiter.TenantName,
					Type:            "waitlist_offered",
					StatusCode:      "waitlist_offered",
					StatusLabel:     "ถึงคิวของคุณแล้ว",
					Title:           "ถึงคิวแผง " + cleanIdent + " แล้ว",
					Message:         fmt.Sprintf("แผง %s (%s) ที่คุณลงชื่อรอคิวว่างลงแล้ว กรุณากดเพื่อยื่นคำร้องขอดูแผงภายใน %d ชั่วโมง (หมดเขต %s น.)", cleanIdent, mktLabel, offerHours, expires.Format("02/01/2006 15:04")),
					StallNumber:     cleanIdent,
					MarketName:      nextWaiter.MarketName,
					ProjectName:     nextWaiter.ProjectName,
					TargetURL:       "/bookings?tab=interests",
				},
				{
					RecipientRole: "officer",
					RecipientName: "เจ้าหน้าที่ดูแลตลาด",
					Type:          "waitlist_offered",
					StatusCode:    "waitlist_offered",
					StatusLabel:   "กำลังเสนอสิทธิ์คิวรอ",
					Title:         fmt.Sprintf("ระบบเสนอสิทธิ์แผง %s ให้คิวที่ 1 (%s)", cleanIdent, nextWaiter.TenantName),
					Message:       fmt.Sprintf("แผง %s ว่างลง ระบบแจ้งเตือนเสนอสิทธิ์ให้คุณ%s (คิวลำดับที่ 1) ตอบรับภายใน %d ชั่วโมง", cleanIdent, nextWaiter.TenantName, offerHours),
					StallNumber:   cleanIdent,
					MarketName:    nextWaiter.MarketName,
					ProjectName:   nextWaiter.ProjectName,
					TargetURL:     "/bookings?role=officer",
				},
			}
			return tx.Create(&notifs).Error
		}
	}

	// 5. No one in Waitlist -> check legacy queue in models.Booking
	var waitingItems []models.Booking
	if err := tx.Where("stall_ident = ? AND status = 'A' AND (queue_number > 0 OR booking_status = 'waiting_queue') AND booking_status NOT IN ('rejected', 'cancelled', 'not_interested', 'no_show', 'lease_applied', 'completed', 'terminated')", cleanIdent).
		Order("queue_number ASC, id ASC").Find(&waitingItems).Error; err != nil {
		return err
	}
	if len(waitingItems) == 0 {
		return tx.Model(&models.Stall{}).Where("ident = ?", cleanIdent).Update("stall_status", "available").Error
	}

	return nil
}

func (r *repository) ProcessWaitlistLifecycle() error {
	now := time.Now()
	return r.conn.Transaction(func(tx *gorm.DB) error {
		// 0. Reconcile queues left behind by applications submitted before the
		// dedicated Waitlist flow was introduced. Once a primary booking has
		// reached lease_applied/completed, no active queue may remain.
		var staleWaitlist []models.Waitlist
		if err := tx.Where(`status IN ('waiting', 'offered') AND (
			EXISTS (
				SELECT 1 FROM booking
				WHERE booking.stall_ident = waitlist.stall_ident
				AND booking.status = 'A'
				AND booking.queue_number = 0
				AND booking.booking_status IN ('lease_applied', 'completed')
			) OR EXISTS (
				SELECT 1 FROM stall
				WHERE stall.ident = waitlist.stall_ident
				AND stall.stall_status = 'occupied'
			)
		)`).Order("stall_ident ASC, position ASC, id ASC").Find(&staleWaitlist).Error; err != nil {
			return err
		}
		for i := range staleWaitlist {
			item := &staleWaitlist[i]
			if err := tx.Model(item).Updates(map[string]interface{}{
				"status":         "cancelled",
				"dismiss_reason": "ตัดสิทธิ์เนื่องจากผู้จองหลักทำสัญญาเช่า/ยื่นคำร้องสำเร็จ",
				"expires_at":     nil,
				"updated_at":     now,
			}).Error; err != nil {
				return err
			}
			notif := models.Notification{
				RecipientRole:   "tenant",
				RecipientUserID: &item.TenantID,
				RecipientName:   item.TenantName,
				Type:            "waitlist_cancelled",
				StatusCode:      "waitlist_cancelled",
				StatusLabel:     "แผงถูกเช่าแล้ว",
				Title:           "แผง " + item.StallIdent + " ถูกเช่าแล้ว",
				Message:         fmt.Sprintf("แผง %s ที่คุณลงชื่อรอคิวมีผู้จองหลักยื่นคำร้องขอเช่าแล้ว ระบบจึงปิดรายการรอคิวของคุณ คุณสามารถเลือกแผงอื่นที่ว่างหรือเข้าคิวแผงอื่นได้", item.StallIdent),
				StallNumber:     item.StallIdent,
				MarketName:      item.MarketName,
				ProjectName:     item.ProjectName,
				TargetURL:       "/bookings?tab=interests",
			}
			if err := tx.Create(&notif).Error; err != nil {
				return err
			}
		}

		// 1. Expire offers where expires_at <= now
		var expiredOffers []models.Waitlist
		if err := tx.Where("status = 'offered' AND expires_at IS NOT NULL AND expires_at <= ?", now).Find(&expiredOffers).Error; err != nil {
			return err
		}
		for i := range expiredOffers {
			item := &expiredOffers[i]
			if err := tx.Model(item).Updates(map[string]interface{}{
				"status":         "expired",
				"dismiss_reason": "หมดเวลาตอบรับสิทธิ์ 48 ชั่วโมง (สละสิทธิ์)",
				"updated_at":     now,
			}).Error; err != nil {
				return err
			}
			expiredNotif := models.Notification{
				RecipientRole:   "tenant",
				RecipientUserID: &item.TenantID,
				RecipientName:   item.TenantName,
				Type:            "waitlist_expired",
				StatusCode:      "waitlist_expired",
				StatusLabel:     "หมดสิทธิ์คิวรอ",
				Title:           "หมดเวลาตอบรับสิทธิ์แผง " + item.StallIdent,
				Message:         fmt.Sprintf("สิทธิ์การยื่นคำร้องขอดูแผง %s ของคุณหมดเวลาแล้ว ระบบได้เลื่อนสิทธิ์ไปยังผู้รอคิวลำดับถัดไป", item.StallIdent),
				StallNumber:     item.StallIdent,
				MarketName:      item.MarketName,
				ProjectName:     item.ProjectName,
				TargetURL:       "/bookings?tab=interests",
			}
			_ = tx.Create(&expiredNotif).Error
			if err := promoteOrReleaseStallOnTx(tx, item.StallIdent); err != nil {
				return err
			}
		}

		// 2. Near-expiry notifications (< 12 hours remaining)
		nearExpiryCutoff := now.Add(12 * time.Hour)
		var nearExpiryOffers []models.Waitlist
		if err := tx.Where("status = 'offered' AND near_expiry_notified = ? AND expires_at IS NOT NULL AND expires_at > ? AND expires_at <= ?", false, now, nearExpiryCutoff).
			Find(&nearExpiryOffers).Error; err != nil {
			return err
		}
		for i := range nearExpiryOffers {
			item := &nearExpiryOffers[i]
			if err := tx.Model(item).Update("near_expiry_notified", true).Error; err != nil {
				return err
			}
			remHours := int(item.ExpiresAt.Sub(now).Hours())
			if remHours < 1 {
				remHours = 1
			}
			warnNotif := models.Notification{
				RecipientRole:   "tenant",
				RecipientUserID: &item.TenantID,
				RecipientName:   item.TenantName,
				Type:            "waitlist_expiring_soon",
				StatusCode:      "waitlist_expiring_soon",
				StatusLabel:     "ใกล้หมดเวลาตอบรับ",
				Title:           "ใกล้หมดเวลายืนยันสิทธิ์คิวแผง " + item.StallIdent,
				Message:         fmt.Sprintf("เหลือเวลาอีกประมาณ %d ชั่วโมงในการยื่นคำร้องขอดูแผง %s กรุณาดำเนินการก่อนหมดสิทธิ์", remHours, item.StallIdent),
				StallNumber:     item.StallIdent,
				MarketName:      item.MarketName,
				ProjectName:     item.ProjectName,
				TargetURL:       "/bookings?tab=interests",
			}
			_ = tx.Create(&warnNotif).Error
		}

		// 3. Auto-expire 'interested' bookings whose lease application deadline (expires_at) has passed
		var expiredInterested []models.Booking
		if err := tx.Where("status = 'A' AND booking_status = 'interested' AND expires_at IS NOT NULL AND expires_at <= ?", now).
			Find(&expiredInterested).Error; err != nil {
			return err
		}
		for i := range expiredInterested {
			bk := &expiredInterested[i]
			autoReason := "เกินกำหนดเวลายื่นคำร้องขอเช่าแผงค้า ระบบตัดสิทธิ์และเลื่อนคิวอัตโนมัติ"
			var tl []responses.TimelineItem
			if bk.Timeline != "" {
				_ = json.Unmarshal([]byte(bk.Timeline), &tl)
			}
			tl = append(tl, responses.TimelineItem{
				Action:    fmt.Sprintf("ตัดสิทธิ์อัตโนมัติแผง %s (เกินกำหนดเวลายื่นคำร้องขอเช่า)", bk.StallIdent),
				ActorRole: "system",
				ActorName: "ระบบจัดการคิวอัตโนมัติ",
				Timestamp: now.Format("2006-01-02T15:04:05Z07:00"),
				Note:      autoReason,
			})
			tlBytes, _ := json.Marshal(tl)
			if err := tx.Model(bk).Updates(map[string]interface{}{
				"booking_status":   "cancelled",
				"rejection_reason": autoReason,
				"note":             autoReason,
				"timeline":         string(tlBytes),
				"updated_at":       now,
			}).Error; err != nil {
				return err
			}
			bkID := bk.ID
			expNotif := models.Notification{
				RecipientRole:   "tenant",
				RecipientUserID: bk.CreatedBy,
				RecipientName:   bk.ApplicantName,
				Type:            "cancelled",
				StatusCode:      "cancelled",
				StatusLabel:     "หมดเวลายื่นคำร้อง",
				Title:           fmt.Sprintf("หมดเวลายื่นคำร้องขอเช่าแผง %s (%s)", bk.StallIdent, bk.BookingNo),
				Message:         fmt.Sprintf("รายการจองแผง %s สิ้นสุดลงเนื่องจากเกินกำหนดเวลายื่นคำร้องขอเช่า ระบบได้ปล่อยสิทธิ์ให้คิวถัดไปโดยอัตโนมัติ", bk.StallIdent),
				BookingID:       &bkID,
				BookingNumber:   bk.BookingNo,
				StallNumber:     bk.StallIdent,
				MarketName:      bk.MarketName,
				ProjectName:     bk.ProjectName,
				TargetURL:       "/bookings/" + fmt.Sprint(bk.ID) + "?role=tenant",
			}
			_ = tx.Create(&expNotif).Error
			if err := promoteOrReleaseStallOnTx(tx, bk.StallIdent); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) GetPolicySettings(marketID int, marketName string) (leaseDays int, offerHours int, resolvedMarketID int, resolvedMarketName string) {
	leaseDays = 3
	offerHours = 48
	resolvedMarketID = marketID
	resolvedMarketName = strings.TrimSpace(marketName)

	var mkt models.Market
	query := r.conn.Model(&models.Market{})
	if marketID > 0 {
		query = query.Where("id = ?", marketID)
	} else if resolvedMarketName != "" {
		query = query.Where("name = ? OR slug = ?", resolvedMarketName, resolvedMarketName)
	} else {
		query = query.Order("id ASC")
	}
	if err := query.First(&mkt).Error; err == nil {
		resolvedMarketID = mkt.ID
		resolvedMarketName = mkt.Name
		if mkt.LeaseApplicationDeadlineDays > 0 {
			leaseDays = mkt.LeaseApplicationDeadlineDays
		}
		if mkt.WaitlistOfferHours > 0 {
			offerHours = mkt.WaitlistOfferHours
		}
	}
	if resolvedMarketName == "" {
		resolvedMarketName = "ทุกตลาดในระบบ"
	}
	return
}

func (r *repository) UpdatePolicySettings(marketID int, marketName string, leaseDays int, offerHours int) (*responses.BookingPolicySettingsResponse, error) {
	if leaseDays <= 0 {
		leaseDays = 3
	}
	if offerHours <= 0 {
		offerHours = 48
	}
	cleanName := strings.TrimSpace(marketName)
	now := time.Now()
	err := r.conn.Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&models.Market{})
		if marketID > 0 {
			query = query.Where("id = ?", marketID)
		} else if cleanName != "" && cleanName != "all" {
			query = query.Where("name = ? OR slug = ?", cleanName, cleanName)
		} else {
			query = query.Where("1 = 1")
		}
		if err := query.Updates(map[string]interface{}{
			"lease_application_deadline_days": leaseDays,
			"waitlist_offer_hours":            offerHours,
			"updated_at":                      now,
		}).Error; err != nil {
			return err
		}

		// Also update active 'interested' bookings so their deadline immediately reflects the updated policy
		bkQuery := tx.Model(&models.Booking{}).Where("status = 'A' AND booking_status = 'interested'")
		if marketID > 0 {
			bkQuery = bkQuery.Where("market_id = ?", marketID)
		} else if cleanName != "" && cleanName != "all" {
			bkQuery = bkQuery.Where("market_name = ?", cleanName)
		}
		var activeInterested []models.Booking
		if err := bkQuery.Find(&activeInterested).Error; err == nil {
			for i := range activeInterested {
				baseTime := activeInterested[i].UpdatedAt
				if baseTime.IsZero() {
					baseTime = now
				}
				newExp := baseTime.Add(time.Duration(leaseDays) * 24 * time.Hour)
				if newExp.Before(now) {
					newExp = now.Add(time.Duration(leaseDays) * 24 * time.Hour)
				}
				_ = tx.Model(&activeInterested[i]).Update("expires_at", newExp).Error
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	lDays, oHours, resID, resName := r.GetPolicySettings(marketID, cleanName)
	return &responses.BookingPolicySettingsResponse{
		MarketID:                     resID,
		MarketName:                   resName,
		LeaseApplicationDeadlineDays: lDays,
		WaitlistOfferHours:           oHours,
	}, nil
}

func (r *repository) FindAll(filter requests.BookingFilter, ownerID *int) ([]models.Booking, error) {
	_ = r.ProcessWaitlistLifecycle()
	var bookings []models.Booking
	query := r.conn.Model(&models.Booking{}).Where("status = 'A'")
	if filter.MarketID != 0 {
		query = query.Where("market_id = ?", filter.MarketID)
	}

	if filter.StallIdent != "" {
		query = query.Where("stall_ident = ?", filter.StallIdent)
	}
	if ownerID != nil {
		query = query.Where("created_by = ?", *ownerID)
	}
	if filter.Status != "" {
		query = query.Where("booking_status = ?", filter.Status)
	}
	if filter.Search != "" {
		query = query.Where("applicant_name LIKE ? OR applicant_phone LIKE ? OR booking_no LIKE ?",
			"%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	err := query.Order("id DESC").Find(&bookings).Error
	return bookings, err
}

func (r *repository) FindByID(id int) (*models.Booking, error) {
	_ = r.ProcessWaitlistLifecycle()
	var booking models.Booking
	err := r.conn.Where("id = ? AND status = 'A'", id).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *repository) FindByGroupID(groupID string) ([]models.Booking, error) {
	var bookings []models.Booking
	clean := strings.TrimSpace(groupID)
	if clean == "" {
		return bookings, nil
	}
	err := r.conn.Where("appointment_group_id = ? AND status = 'A'", clean).Order("id ASC").Find(&bookings).Error
	return bookings, err
}

func (r *repository) FindActiveDuplicate(stallIdent string, userID int, applicantName string, applicantPhone string) (*models.Booking, error) {
	var booking models.Booking
	query := r.conn.Model(&models.Booking{}).
		Where("stall_ident = ? AND status = 'A' AND booking_status NOT IN ('rejected', 'cancelled', 'not_interested', 'no_show', 'terminated')", stallIdent)

	if userID > 0 {
		query = query.Where("created_by = ?", userID)
	} else {
		query = query.Where("(applicant_phone = ? OR applicant_name = ?)", applicantPhone, applicantName)
	}

	err := query.Order("id DESC").First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *repository) HasActiveHolder(stallIdent string) (bool, error) {
	var count int64
	err := r.conn.Model(&models.Booking{}).
		Where("stall_ident = ? AND status = 'A' AND queue_number = 0 AND booking_status NOT IN ('waiting_queue', 'rejected', 'cancelled', 'not_interested', 'no_show', 'terminated')", stallIdent).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) LookupUserByID(userID int) (*models.User, error) {
	if userID <= 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var user models.User
	if err := r.conn.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) LookupStallRate(stallIdent string, marketID *int) (monthlyRent float64, depositFee float64, projectName string, marketName string) {
	monthlyRent = 1800
	depositFee = 1800
	var stall models.Stall
	query := r.conn.Where("ident = ? AND status = 'A'", strings.TrimSpace(stallIdent))
	if marketID != nil && *marketID > 0 {
		query = query.Where("market_id = ?", *marketID)
	}
	if err := query.First(&stall).Error; err == nil {
		if stall.MonthlyRent > 0 {
			monthlyRent = stall.MonthlyRent
			depositFee = stall.MonthlyRent
		}
		projectName = stall.ProjectName
		marketName = stall.MarketName
	}
	return
}

func (r *repository) SyncBookingStalls(parentBookingID int, parentBookingNo string, groupID string, bookings []models.Booking) {
	if len(bookings) == 0 {
		return
	}
	for _, b := range bookings {
		rent := b.MonthlyRent
		if rent <= 0 {
			rent, _, _, _ = r.LookupStallRate(b.StallIdent, b.MarketID)
		}
		dep := b.DepositFee
		if dep <= 0 {
			dep = rent
		}
		mktID := uint(1)
		if b.MarketID != nil && *b.MarketID > 0 {
			mktID = uint(*b.MarketID)
		}
		stID := uint(0)
		if b.StallID != nil && *b.StallID > 0 {
			stID = uint(*b.StallID)
		}
		var existing models.BookingStall
		err := r.conn.Where("sub_booking_id = ?", b.ID).First(&existing).Error
		var cancelledAt *time.Time
		if b.BookingStatus == "cancelled" || b.BookingStatus == "rejected" || b.BookingStatus == "not_interested" || b.BookingStatus == "no_show" {
			t := b.UpdatedAt
			cancelledAt = &t
		}
		record := models.BookingStall{
			ID:                 existing.ID,
			BookingID:          uint(parentBookingID),
			SubBookingID:       uint(b.ID),
			BookingNo:          parentBookingNo,
			AppointmentGroupID: groupID,
			MarketID:           mktID,
			MarketName:         b.MarketName,
			ProjectName:        b.ProjectName,
			StallID:            stID,
			StallIdent:         b.StallIdent,
			MonthlyRent:        rent,
			DepositFee:         dep,
			Status:             b.BookingStatus,
			CancelReason:       b.RejectionReason,
			CancelledAt:        cancelledAt,
			CreatedAt:          b.CreatedAt,
			UpdatedAt:          b.UpdatedAt,
		}
		if err == nil && existing.ID > 0 {
			_ = r.conn.Save(&record).Error
		} else {
			_ = r.conn.Create(&record).Error
		}
		if b.ProjectName != "" {
			_ = r.conn.Model(&models.Booking{}).Where("id = ? AND project_name <> ?", b.ID, b.ProjectName).Update("project_name", b.ProjectName).Error
			_ = r.conn.Model(&models.Stall{}).Where("ident = ? AND project_name <> ?", b.StallIdent, b.ProjectName).Update("project_name", b.ProjectName).Error
		}
	}
}

func (r *repository) Create(booking *models.Booking) error {
	now := time.Now()
	booking.CreatedAt = now
	booking.UpdatedAt = now
	booking.Status = "A"
	if booking.MonthlyRent <= 0 {
		rent, dep, projName, mktName := r.LookupStallRate(booking.StallIdent, booking.MarketID)
		booking.MonthlyRent = rent
		booking.DepositFee = dep
		if booking.ProjectName == "" && projName != "" {
			booking.ProjectName = projName
		}
		if booking.MarketName == "" && mktName != "" {
			booking.MarketName = mktName
		}
	} else if booking.DepositFee <= 0 {
		booking.DepositFee = booking.MonthlyRent
	}
	return r.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		if !booking.SkipNotification {
			bookingID := booking.ID
			statusLabel := "รออนุมัติวันนัด"
			officerTitle := "มีคำขอนัดชมแผงใหม่"
			officerMsg := booking.ApplicantName + " ส่งคำขอนัดชมแผง " + booking.StallIdent + " (ผู้จองหลัก)"

			officerNotif := models.Notification{
				RecipientRole: "officer", RecipientName: "เจ้าหน้าที่ดูแลตลาด",
				Type: booking.BookingStatus, StatusCode: booking.BookingStatus, StatusLabel: statusLabel,
				Title:     officerTitle + " " + booking.StallIdent + " (" + booking.BookingNo + ")",
				Message:   officerMsg,
				BookingID: &bookingID, BookingNumber: booking.BookingNo, StallNumber: booking.StallIdent,
				MarketName: booking.MarketName, ProjectName: booking.ProjectName,
				TargetURL: "/bookings/" + fmt.Sprint(booking.ID) + "?role=officer",
			}
			if err := tx.Create(&officerNotif).Error; err != nil {
				return err
			}
		}
		var stallCount int64
		tx.Model(&models.Stall{}).Where("ident = ?", booking.StallIdent).Count(&stallCount)
		if stallCount == 0 {
			marketID := 1
			if booking.MarketID != nil {
				marketID = *booking.MarketID
			}
			newStall := models.Stall{
				MarketID:    marketID,
				MarketName:  booking.MarketName,
				ProjectName: booking.ProjectName,
				Ident:       booking.StallIdent,
				Name:        "แผงค้า " + booking.StallIdent,
				MonthlyRent: booking.MonthlyRent,
				StallStatus: "reserved",
				Floor:       "1",
				Status:      "A",
			}
			return tx.Create(&newStall).Error
		}
		return tx.Model(&models.Stall{}).
			Where("ident = ?", booking.StallIdent).
			Update("stall_status", "reserved").Error
	})
}

func appendTimelineOnModel(b *models.Booking, item responses.TimelineItem) {
	var tl []responses.TimelineItem
	if b.Timeline != "" {
		_ = json.Unmarshal([]byte(b.Timeline), &tl)
	}
	tl = append(tl, item)
	if bytes, err := json.Marshal(tl); err == nil {
		b.Timeline = string(bytes)
	}
}

func (r *repository) Update(booking *models.Booking) error {
	now := time.Now()
	booking.UpdatedAt = now
	isTerminalCancel := booking.BookingStatus == "cancelled" ||
		booking.BookingStatus == "rejected" ||
		booking.BookingStatus == "not_interested" ||
		booking.BookingStatus == "no_show"
	if isTerminalCancel {
		booking.QueueNumber = 0
		booking.ExpiresAt = nil
	} else if booking.BookingStatus == "interested" {
		mID := 0
		if booking.MarketID != nil {
			mID = *booking.MarketID
		}
		leaseDays, _, _, _ := r.GetPolicySettings(mID, booking.MarketName)
		exp := now.Add(time.Duration(leaseDays) * 24 * time.Hour)
		booking.ExpiresAt = &exp
	} else if booking.BookingStatus == "lease_applied" || booking.BookingStatus == "completed" {
		booking.ExpiresAt = nil
	}

	return r.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(booking).Error; err != nil {
			return err
		}

		statusLabels := map[string]string{
			"approved":         "นัดหมายสำเร็จ",
			"rejected":         "ไม่อนุมัตินัดหมาย",
			"interested":       "พร้อมยื่นคำร้องขอเช่า",
			"cancelled":        "ยกเลิกรายการแล้ว",
			"pending_approval": "รออนุมัติวันนัด",
			"no_show":          "ไม่มาตามนัดหมาย",
			"not_interested":   "สละสิทธิ์การเช่า",
		}
		label := statusLabels[booking.BookingStatus]
		if label == "" {
			label = booking.BookingStatus
		}
		notifTitle := label + " " + booking.StallIdent + " (" + booking.BookingNo + ")"
		notifMsg := "สถานะรายการแผง " + booking.StallIdent + " เปลี่ยนเป็น " + label
		if booking.BookingStatus == "interested" {
			mID := 0
			if booking.MarketID != nil {
				mID = *booking.MarketID
			}
			leaseDays, _, _, _ := r.GetPolicySettings(mID, booking.MarketName)
			notifTitle = "เปิดสิทธิ์ยื่นคำร้องขอเช่าแผง " + booking.StallIdent + " (" + booking.BookingNo + ")"
			notifMsg = fmt.Sprintf("เจ้าหน้าที่ยืนยันการเช่าแผง %s เรียบร้อยแล้ว กรุณากดปุ่ม 'ยื่นคำร้องขอเช่าแผงนี้' ภายใน %d วัน", booking.StallIdent, leaseDays)
		}

		if !booking.SkipNotification {
			bookingID := booking.ID
			notification := models.Notification{
				RecipientRole: "tenant", RecipientUserID: booking.CreatedBy, RecipientName: booking.ApplicantName,
				Type: booking.BookingStatus, StatusCode: booking.BookingStatus, StatusLabel: label,
				Title:     notifTitle,
				Message:   notifMsg,
				BookingID: &bookingID, BookingNumber: booking.BookingNo, StallNumber: booking.StallIdent,
				MarketName: booking.MarketName, ProjectName: booking.ProjectName,
				TargetURL: func() string {
					if booking.BookingStatus == "interested" {
						return "/applications?booking=" + fmt.Sprint(booking.ID)
					}
					return "/bookings/" + fmt.Sprint(booking.ID) + "?role=tenant"
				}(),
			}
			if err := tx.Create(&notification).Error; err != nil {
				return err
			}
		}

		if booking.BookingStatus == "completed" || booking.BookingStatus == "lease_applied" {
			return tx.Model(&models.Stall{}).Where("ident = ?", booking.StallIdent).Update("stall_status", "occupied").Error
		}

		if !isTerminalCancel {
			return nil
		}

		// When a stall is cancelled/released by its holder, offer to #1 in Waitlist or return stall to available
		return promoteOrReleaseStallOnTx(tx, booking.StallIdent)
	})
}

func (r *repository) CountActiveQueue(stallIdent string) (int64, error) {
	return r.CountWaitlistByStall(stallIdent)
}

func (r *repository) CountWaitlistByStall(stallIdent string) (int64, error) {
	var count int64
	err := r.conn.Model(&models.Waitlist{}).
		Where("stall_ident = ? AND status IN ('waiting', 'offered')", strings.TrimSpace(stallIdent)).
		Count(&count).Error
	return count, err
}

func (r *repository) UpdateStallStatus(stallIdent string, status string) error {
	return r.conn.Model(&models.Stall{}).
		Where("ident = ?", stallIdent).
		Update("stall_status", status).Error
}

func (r *repository) ClearAll() error {
	return r.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("1 = 1").Delete(&models.Notification{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1 = 1").Delete(&models.Contract{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1 = 1").Delete(&models.Application{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1 = 1").Delete(&models.BookingStall{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1 = 1").Delete(&models.Waitlist{}).Error; err != nil {
			return err
		}
		if err := tx.Where("1 = 1").Delete(&models.Booking{}).Error; err != nil {
			return err
		}
		return tx.Model(&models.Stall{}).Where("stall_status IN ('reserved', 'booked', 'occupied', 'offered')").Update("stall_status", "available").Error
	})
}

func (r *repository) CreateWaitlist(item *models.Waitlist) error {
	_ = r.ProcessWaitlistLifecycle()
	return r.conn.Transaction(func(tx *gorm.DB) error {
		var existing int64
		if err := tx.Model(&models.Waitlist{}).
			Where("tenant_id = ? AND stall_ident = ? AND status IN ('waiting', 'offered')", item.TenantID, item.StallIdent).
			Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			return fmt.Errorf("WAITLIST_DUPLICATE")
		}
		var last models.Waitlist
		position := 1
		if err := tx.Where("stall_ident = ? AND status = 'waiting'", item.StallIdent).
			Order("position DESC, id DESC").First(&last).Error; err == nil {
			position = last.Position + 1
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		item.Position = position
		item.Status = "waiting"
		return tx.Create(item).Error
	})
}

func (r *repository) FindWaitlistByTenant(userID int) ([]models.Waitlist, error) {
	_ = r.ProcessWaitlistLifecycle()
	var items []models.Waitlist
	query := r.conn.Where("1 = 1")
	if userID > 0 {
		query = query.Where("tenant_id = ? AND tenant_dismissed = false AND status NOT IN ('withdrawn', 'accepted')", userID)
	} else {
		query = query.Where("status IN ('waiting', 'offered', 'cancelled', 'expired')")
	}
	err := query.
		Order("CASE WHEN status = 'offered' THEN 0 WHEN status = 'waiting' THEN 1 ELSE 2 END ASC, position ASC, updated_at DESC, created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *repository) FindWaitlistByStall(stallIdent string) ([]models.Waitlist, error) {
	_ = r.ProcessWaitlistLifecycle()
	var items []models.Waitlist
	err := r.conn.Where("stall_ident = ?", strings.TrimSpace(stallIdent)).
		Order("CASE WHEN status = 'offered' THEN 0 WHEN status = 'waiting' THEN 1 ELSE 2 END ASC, position ASC, updated_at DESC, id ASC").
		Find(&items).Error
	return items, err
}

func (r *repository) FindWaitlistByID(id int) (*models.Waitlist, error) {
	_ = r.ProcessWaitlistLifecycle()
	var item models.Waitlist
	if err := r.conn.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) WithdrawWaitlist(item *models.Waitlist) error {
	return r.conn.Transaction(func(tx *gorm.DB) error {
		wasOffered := item.Status == "offered"
		wasActive := item.Status == "waiting" || item.Status == "offered"
		now := time.Now()

		updates := map[string]interface{}{
			"tenant_dismissed": true,
			"updated_at":       now,
		}
		if wasActive {
			dismissReason := item.DismissReason
			if dismissReason == "" {
				dismissReason = "ผู้เช่ายกเลิกการรอคิวด้วยตนเอง"
			}
			updates["status"] = "withdrawn"
			updates["dismiss_reason"] = dismissReason
		}

		if err := tx.Model(item).Updates(updates).Error; err != nil {
			return err
		}

		// Re-index remaining active waitlist items only if the withdrawing item was active
		if wasActive {
			// If the person who withdrew was holding the active offer, promote the next person or release the stall!
			if wasOffered {
				return promoteOrReleaseStallOnTx(tx, item.StallIdent)
			}

			// If a waiting person withdrew, re-index remaining waiting items
			var remainingWaiters []models.Waitlist
			if err := tx.Where("stall_ident = ? AND status = 'waiting'", item.StallIdent).
				Order("position ASC, id ASC").Find(&remainingWaiters).Error; err != nil {
				return err
			}
			for i := range remainingWaiters {
				newPos := i + 1
				if remainingWaiters[i].Position != newPos {
					if err := tx.Model(&remainingWaiters[i]).Updates(map[string]interface{}{
						"position":   newPos,
						"updated_at": now,
					}).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func (r *repository) AdvanceWaitlistByStall(stallIdent string) error {
	cleanIdent := strings.TrimSpace(stallIdent)
	now := time.Now()
	return r.conn.Transaction(func(tx *gorm.DB) error {
		var currentOffers []models.Waitlist
		if err := tx.Where("stall_ident = ? AND status = 'offered'", cleanIdent).Find(&currentOffers).Error; err != nil {
			return err
		}
		for i := range currentOffers {
			offer := &currentOffers[i]
			if err := tx.Model(offer).Updates(map[string]interface{}{
				"status":         "expired",
				"dismiss_reason": "เจ้าหน้าที่ข้ามสิทธิ์และเลื่อนให้คิวถัดไป",
				"updated_at":     now,
			}).Error; err != nil {
				return err
			}
			notif := models.Notification{
				RecipientRole:   "tenant",
				RecipientUserID: &offer.TenantID,
				RecipientName:   offer.TenantName,
				Type:            "waitlist_expired",
				StatusCode:      "waitlist_expired",
				StatusLabel:     "หมดสิทธิ์คิวรอ",
				Title:           "ตัดสิทธิ์/เลื่อนคิวแผง " + cleanIdent,
				Message:         fmt.Sprintf("สิทธิ์คิวแผง %s ของคุณสิ้นสุดลงเนื่องจากหมดเวลาหรือสละสิทธิ์ ระบบได้เลื่อนสิทธิ์ให้คิวถัดไปแล้ว", cleanIdent),
				StallNumber:     cleanIdent,
				MarketName:      offer.MarketName,
				ProjectName:     offer.ProjectName,
				TargetURL:       "/bookings?tab=interests",
			}
			_ = tx.Create(&notif).Error
		}
		return promoteOrReleaseStallOnTx(tx, cleanIdent)
	})
}

func (r *repository) FindActiveWaitlistOffer(stallIdent string) (*models.Waitlist, error) {
	_ = r.ProcessWaitlistLifecycle()
	var item models.Waitlist
	err := r.conn.Where("stall_ident = ? AND status = 'offered' AND expires_at > ?", strings.TrimSpace(stallIdent), time.Now()).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) AcceptWaitlistOffer(stallIdent string, tenantID int) error {
	now := time.Now()
	return r.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Waitlist{}).
			Where("stall_ident = ? AND tenant_id = ? AND status = 'offered'", strings.TrimSpace(stallIdent), tenantID).
			Updates(map[string]interface{}{"status": "accepted", "updated_at": now}).Error; err != nil {
			return err
		}
		// Re-index any remaining waiting items on this stall (1, 2, ...)
		var remaining []models.Waitlist
		if err := tx.Where("stall_ident = ? AND status = 'waiting'", strings.TrimSpace(stallIdent)).
			Order("position ASC, id ASC").Find(&remaining).Error; err == nil {
			for i := range remaining {
				_ = tx.Model(&remaining[i]).Update("position", i+1).Error
			}
		}
		return nil
	})
}

func (r *repository) CreateNotification(notif *models.Notification) error {
	return r.conn.Create(notif).Error
}

func (r *repository) CreateNotifications(notifs []models.Notification) error {
	return r.conn.Create(&notifs).Error
}
