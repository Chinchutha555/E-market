package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/application/domains"
	"gorm.io/gorm"
)

type repository struct {
	conn *gorm.DB
}

func NewRepository(db *gorm.DB) domains.Repository {
	return &repository{conn: db}
}

// fillDocumentNumbers assigns a yearly sequence based on application creation order.
func (r *repository) fillDocumentNumbers(apps []models.Application) error {
	if len(apps) == 0 {
		return nil
	}
	var all []models.Application
	if err := r.conn.Select("id", "created_at").Order("id ASC").Find(&all).Error; err != nil {
		return err
	}
	zone := time.FixedZone("Asia/Bangkok", 7*60*60)
	counts := make(map[int]int)
	numbers := make(map[int]string, len(all))
	for _, app := range all {
		year := app.CreatedAt.In(zone).Year()
		counts[year]++
		numbers[app.ID] = fmt.Sprintf("%d/%02d", counts[year], (year+543)%100)
	}
	for i := range apps {
		apps[i].DocumentNumber = numbers[apps[i].ID]
	}
	return nil
}

// fillRentalTypes resolves the officer's selection for legacy monthly applications.
func (r *repository) fillRentalTypes(apps []models.Application) error {
	ids := make([]int, 0, len(apps))
	for _, app := range apps {
		if app.BookingID != nil {
			ids = append(ids, *app.BookingID)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	var bookings []models.Booking
	if err := r.conn.Select("id", "rental_type", "booking_no").Where("id IN ?", ids).Find(&bookings).Error; err != nil {
		return err
	}
	numbers := make(map[int]string, len(bookings))
	types := make(map[int]string, len(bookings))
	for _, booking := range bookings {
		numbers[booking.ID] = booking.BookingNo
		if booking.RentalType == "trial" || booking.RentalType == "temporary" || booking.RentalType == "contract" {
			types[booking.ID] = booking.RentalType
		}
	}
	for i := range apps {
		if apps[i].BookingID != nil {
			apps[i].BookingNumber = numbers[*apps[i].BookingID]
		}
		if apps[i].BookingID != nil && (apps[i].ProposedRentType == "" || apps[i].ProposedRentType == "monthly") {
			if value := types[*apps[i].BookingID]; value != "" {
				apps[i].ProposedRentType = value
			}
		}
	}
	return nil
}

func (r *repository) FindAll(filter requests.ApplicationFilter) ([]models.Application, error) {
	var apps []models.Application
	query := r.conn.Model(&models.Application{}).Where("(status = 'A' OR status IS NULL OR status = '')")
	if filter.MarketID != 0 {
		query = query.Where("market_id = ?", filter.MarketID)
	}
	if filter.CreatedBy != 0 {
		query = query.Where("created_by = ?", filter.CreatedBy)
	}

	if filter.StallIdent != "" {
		query = query.Where("stall_ident = ?", filter.StallIdent)
	}
	if filter.Status != "" {
		query = query.Where("application_state = ?", filter.Status)
	}
	if filter.ApplicationType != "" {
		if filter.ApplicationType == "general" {
			query = query.Where("application_type = 'general' OR application_no LIKE 'GEN-%'")
		} else if filter.ApplicationType == "stall_lease" {
			query = query.Where("(application_type = 'stall_lease' OR application_type = '' OR application_type IS NULL) AND application_no NOT LIKE 'GEN-%'")
		} else {
			query = query.Where("application_type = ?", filter.ApplicationType)
		}
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("LOWER(applicant_name) LIKE LOWER(?) OR applicant_phone LIKE ? OR LOWER(shop_name) LIKE LOWER(?) OR LOWER(application_no) LIKE LOWER(?) OR LOWER(subject) LIKE LOWER(?)",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	err := query.Order("id DESC").Find(&apps).Error
	if err == nil {
		for i := range apps {
			if apps[i].MarketName == "" && apps[i].MarketID != nil && *apps[i].MarketID > 0 {
				var mkt models.Market
				if r.conn.Where("id = ?", *apps[i].MarketID).First(&mkt).Error == nil {
					apps[i].MarketName = mkt.Name
				}
			}
		}
	}
	if err == nil {
		err = r.fillRentalTypes(apps)
	}
	if err == nil {
		err = r.fillDocumentNumbers(apps)
	}
	return apps, err
}

func (r *repository) FindByID(id int) (*models.Application, error) {
	var app models.Application
	err := r.conn.Where("id = ? AND (status = 'A' OR status IS NULL OR status = '')", id).First(&app).Error
	if err != nil {
		return nil, err
	}
	if app.MarketName == "" && app.MarketID != nil && *app.MarketID > 0 {
		var mkt models.Market
		if r.conn.Where("id = ?", *app.MarketID).First(&mkt).Error == nil {
			app.MarketName = mkt.Name
		}
	}
	items := []models.Application{app}
	if err := r.fillRentalTypes(items); err != nil {
		return nil, err
	}
	if err := r.fillDocumentNumbers(items); err != nil {
		return nil, err
	}
	app.DocumentNumber = items[0].DocumentNumber
	app.ProposedRentType = items[0].ProposedRentType
	app.BookingNumber = items[0].BookingNumber
	return &app, nil
}

func (r *repository) Create(app *models.Application) error {
	now := time.Now()
	app.CreatedAt = now
	app.UpdatedAt = now
	app.Status = "A"

	// Resolve MarketID and MarketName
	if app.MarketID == nil && app.MarketName != "" {
		var mkt models.Market
		if r.conn.Where("name = ? OR name_en = ?", app.MarketName, app.MarketName).First(&mkt).Error == nil {
			app.MarketID = &mkt.ID
		}
	}
	if app.MarketID == nil {
		var stall models.Stall
		query := r.conn.Where("status = 'A'")
		if app.StallID != nil && *app.StallID > 0 {
			query = query.Where("id = ?", *app.StallID)
		} else if app.StallIdent != "" && app.StallIdent != "-" {
			query = query.Where("ident = ?", app.StallIdent)
		}
		if query.First(&stall).Error == nil && stall.MarketID > 0 {
			app.MarketID = &stall.MarketID
		}
	}
	if app.MarketName == "" && app.MarketID != nil && *app.MarketID > 0 {
		var mkt models.Market
		if r.conn.Where("id = ?", *app.MarketID).First(&mkt).Error == nil {
			app.MarketName = mkt.Name
		}
	}

	return r.conn.Transaction(func(tx *gorm.DB) error {
		if app.BookingID != nil && (app.ProposedRentType == "" || app.ProposedRentType == "monthly") {
			var booking models.Booking
			if err := tx.Select("id", "rental_type").Where("id = ?", *app.BookingID).First(&booking).Error; err == nil {
				if booking.RentalType == "trial" || booking.RentalType == "temporary" || booking.RentalType == "contract" {
					app.ProposedRentType = booking.RentalType
				}
			} else if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		if err := tx.Create(app).Error; err != nil {
			return err
		}

		// Only link bookings and update stall status for stall_lease applications
		if (app.ApplicationType == "stall_lease" || app.ApplicationType == "") && app.StallIdent != "" && app.StallIdent != "-" {
			// 1. Mark the primary holder booking for this stall as lease_applied and append timeline
			var activeBookings []models.Booking
			if app.BookingID != nil && *app.BookingID > 0 {
				tx.Model(&models.Booking{}).
					Where("id = ? OR (stall_ident = ? AND status = 'A' AND queue_number = 0 AND booking_status IN ('confirmed_rent', 'interested', 'approved', 'pending_decision', 'pending_approval'))", *app.BookingID, app.StallIdent).
					Find(&activeBookings)
			} else {
				tx.Model(&models.Booking{}).
					Where("stall_ident = ? AND status = 'A' AND queue_number = 0 AND booking_status IN ('confirmed_rent', 'interested', 'approved', 'pending_decision', 'pending_approval')", app.StallIdent).
					Find(&activeBookings)
			}

			for _, bk := range activeBookings {
				var tl []responses.TimelineItem
				if bk.Timeline != "" {
					_ = json.Unmarshal([]byte(bk.Timeline), &tl)
				}
				actor := app.ApplicantName
				if actor == "" {
					actor = bk.ApplicantName
				}
				if actor == "" {
					actor = "ผู้เช่า"
				}
				tl = append(tl, responses.TimelineItem{
					Action:    "ยื่นคำร้องขอเช่าแผงค้า",
					ActorRole: "tenant",
					ActorName: actor,
					Timestamp: now.Format("2006-01-02T15:04:05Z07:00"),
					Note:      fmt.Sprintf("ส่งคำร้องขอเช่าแผงค้าเข้าสู่ระบบเรียบร้อยแล้ว รหัสคำร้อง: %s", app.ApplicationNo),
				})
				tlBytes, _ := json.Marshal(tl)
				_ = tx.Model(&models.Booking{}).Where("id = ?", bk.ID).Updates(map[string]interface{}{
					"booking_status": "lease_applied",
					"timeline":       string(tlBytes),
					"updated_at":     now,
				}).Error
			}

			// 2. Automatically close all backup waiting queue bookings for this stall and notify each queued tenant
			var backupQueues []models.Booking
			tx.Model(&models.Booking{}).
				Where("stall_ident = ? AND status = 'A' AND (queue_number > 0 OR booking_status = 'waiting_queue') AND booking_status NOT IN ('rejected', 'cancelled', 'not_interested', 'no_show', 'lease_applied', 'completed', 'terminated')", app.StallIdent).
				Find(&backupQueues)

			for _, q := range backupQueues {
				var qTl []responses.TimelineItem
				if q.Timeline != "" {
					_ = json.Unmarshal([]byte(q.Timeline), &qTl)
				}
				qTl = append(qTl, responses.TimelineItem{
					Action:    "ปิดคิวสำรองอัตโนมัติ (ผู้จองหลักยื่นคำร้องเช่าแล้ว)",
					ActorRole: "system",
					ActorName: "ระบบจัดการคิวอัตโนมัติ",
					Timestamp: now.Format("2006-01-02T15:04:05Z07:00"),
					Note:      "ผู้จองคิวหลักได้ยื่นคำร้องขอเช่าแผงนี้สำเร็จแล้ว ระบบจึงปิดคิวสำรองของคุณโดยอัตโนมัติ",
				})
				qTlBytes, _ := json.Marshal(qTl)
				closeNote := "ผู้จองคิวหลักได้ยื่นคำร้องเช่าแผงนี้สำเร็จแล้ว (ปิดคิวสำรองอัตโนมัติ)"
				_ = tx.Model(&models.Booking{}).Where("id = ?", q.ID).Updates(map[string]interface{}{
					"booking_status": "cancelled",
					"queue_number":   0,
					"note":           closeNote,
					"timeline":       string(qTlBytes),
					"updated_at":     now,
				}).Error

				qID := q.ID
				closeNotif := models.Notification{
					RecipientRole:   "tenant",
					RecipientUserID: q.CreatedBy,
					RecipientName:   q.ApplicantName,
					Type:            "cancelled",
					StatusCode:      "cancelled",
					StatusLabel:     "ปิดคิวสำรอง",
					Title:           fmt.Sprintf("ปิดคิวสำรองแผง %s (%s)", q.StallIdent, q.BookingNo),
					Message:         fmt.Sprintf("ผู้จองคิวหลักได้ยื่นคำร้องเช่าแผง %s สำเร็จแล้ว ระบบได้ปิดคิวสำรองของคุณ ท่านสามารถเลือกจองแผงอื่นได้", q.StallIdent),
					BookingID:       &qID,
					BookingNumber:   q.BookingNo,
					StallNumber:     q.StallIdent,
					MarketName:      q.MarketName,
					ProjectName:     q.ProjectName,
					TargetURL:       "/bookings/" + fmt.Sprint(q.ID) + "?role=tenant",
				}
				_ = tx.Create(&closeNotif).Error
			}

			// 3. Close the per-stall Waitlist as soon as the primary holder
			// submits a lease application. These records are separate from the
			// legacy queue bookings above and must be notified individually.
			var activeWaitlist []models.Waitlist
			if err := tx.Where("stall_ident = ? AND status IN ('waiting', 'offered')", app.StallIdent).
				Order("position ASC, id ASC").Find(&activeWaitlist).Error; err != nil {
				return err
			}
			for i := range activeWaitlist {
				item := &activeWaitlist[i]
				if err := tx.Model(item).Updates(map[string]interface{}{
					"status":         "cancelled",
					"dismiss_reason": "ตัดสิทธิ์เนื่องจากผู้จองหลักยื่นคำร้องขอเช่าแผงสำเร็จ",
					"expires_at":     nil,
					"updated_at":     now,
				}).Error; err != nil {
					return err
				}

				closeNotif := models.Notification{
					RecipientRole:   "tenant",
					RecipientUserID: &item.TenantID,
					RecipientName:   item.TenantName,
					Type:            "waitlist_cancelled",
					StatusCode:      "waitlist_cancelled",
					StatusLabel:     "แผงถูกเช่าแล้ว",
					Title:           fmt.Sprintf("แผง %s ถูกเช่าแล้ว", item.StallIdent),
					Message:         fmt.Sprintf("แผง %s ที่คุณลงชื่อรอคิวมีผู้จองหลักยื่นคำร้องขอเช่าแล้ว ระบบจึงปิดรายการรอคิวของคุณ คุณสามารถเลือกแผงอื่นที่ว่างหรือเข้าคิวแผงอื่นได้", item.StallIdent),
					StallNumber:     item.StallIdent,
					MarketName:      item.MarketName,
					ProjectName:     item.ProjectName,
					TargetURL:       "/bookings?tab=interests",
				}
				if err := tx.Create(&closeNotif).Error; err != nil {
					return err
				}
			}

			// 4. Upsert stall to occupied
			var stallCount int64
			tx.Model(&models.Stall{}).Where("ident = ?", app.StallIdent).Count(&stallCount)
			if stallCount == 0 {
				newStall := models.Stall{
					Ident:       app.StallIdent,
					Name:        "แผงค้า " + app.StallIdent,
					StallStatus: "occupied",
					Floor:       "1",
					Status:      "A",
				}
				if app.MarketID != nil {
					newStall.MarketID = *app.MarketID
				}
				_ = tx.Create(&newStall).Error
			} else {
				_ = tx.Model(&models.Stall{}).
					Where("ident = ?", app.StallIdent).
					Update("stall_status", "occupied").Error
			}
		}
		return nil
	})
}

func (r *repository) Update(app *models.Application) error {
	now := time.Now()
	app.UpdatedAt = now
	return r.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(app).Error; err != nil {
			return err
		}
		if app.StallIdent != "" && app.ApplicationState == "rejected" {
			var linkedBookings []models.Booking
			tx.Model(&models.Booking{}).
				Where("stall_ident = ? AND status = 'A' AND booking_status IN ('lease_applied', 'interested')", app.StallIdent).
				Find(&linkedBookings)
			for _, bk := range linkedBookings {
				var tl []responses.TimelineItem
				if bk.Timeline != "" {
					_ = json.Unmarshal([]byte(bk.Timeline), &tl)
				}
				reasonNote := app.ReviewNote
				if reasonNote == "" {
					reasonNote = "คำร้องขอเช่าแผงค้าไม่ผ่านการอนุมัติ"
				}
				tl = append(tl, responses.TimelineItem{
					Action:    "คำร้องขอเช่าไม่ผ่านการอนุมัติ (ปล่อยแผงค้าคืนสู่ระบบ)",
					ActorRole: "officer",
					ActorName: "เจ้าหน้าที่พิจารณาคำร้อง",
					Timestamp: now.Format("2006-01-02T15:04:05Z07:00"),
					Note:      reasonNote,
				})
				tlBytes, _ := json.Marshal(tl)
				_ = tx.Model(&models.Booking{}).Where("id = ?", bk.ID).Updates(map[string]interface{}{
					"booking_status":   "rejected",
					"queue_number":     0,
					"rejection_reason": reasonNote,
					"note":             reasonNote,
					"timeline":         string(tlBytes),
					"updated_at":       now,
				}).Error
			}
		} else if app.StallIdent != "" && app.ApplicationState == "approved" {
			var linkedBookings []models.Booking
			tx.Model(&models.Booking{}).
				Where("stall_ident = ? AND status = 'A' AND booking_status = 'lease_applied'", app.StallIdent).
				Find(&linkedBookings)
			for _, bk := range linkedBookings {
				var tl []responses.TimelineItem
				if bk.Timeline != "" {
					_ = json.Unmarshal([]byte(bk.Timeline), &tl)
				}
				tl = append(tl, responses.TimelineItem{
					Action:    "อนุมัติคำร้องขอเช่าแผงค้าสำเร็จ",
					ActorRole: "officer",
					ActorName: "ผู้อำนวยการสำนักตลาด",
					Timestamp: now.Format("2006-01-02T15:04:05Z07:00"),
					Note:      "คำร้องขอเช่าผ่านการอนุมัติเรียบร้อยแล้ว พร้อมจัดทำสัญญาเช่า",
				})
				tlBytes, _ := json.Marshal(tl)
				_ = tx.Model(&models.Booking{}).Where("id = ?", bk.ID).Updates(map[string]interface{}{
					"booking_status": "completed",
					"timeline":       string(tlBytes),
					"updated_at":     now,
				}).Error
			}
		}
		return nil
	})
}

func (r *repository) UpdateStallStatus(stallIdent string, status string) error {
	var stallCount int64
	r.conn.Model(&models.Stall{}).Where("ident = ?", stallIdent).Count(&stallCount)
	if stallCount == 0 {
		newStall := models.Stall{
			Ident:       stallIdent,
			Name:        "แผงค้า " + stallIdent,
			StallStatus: status,
			Floor:       "1",
			Status:      "A",
		}
		return r.conn.Create(&newStall).Error
	}
	return r.conn.Model(&models.Stall{}).
		Where("ident = ?", stallIdent).
		Update("stall_status", status).Error
}
