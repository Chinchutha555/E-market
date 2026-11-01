package repositories

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/src/stall/domains"
	"gorm.io/gorm"
)

type repository struct {
	conn *gorm.DB
}

func NewRepository(db *gorm.DB) domains.Repository {
	return &repository{conn: db}
}

func (r *repository) FindAll(filter requests.StallFilter) ([]models.Stall, error) {
	var stalls []models.Stall
	query := r.conn.Model(&models.Stall{}).Where("status = 'A'")
	resolvedMarketID := filter.MarketID
	resolvedMarketSlug := ""
	resolvedMarketName := ""

	marketKey := strings.TrimSpace(filter.Market)
	if marketKey == "" {
		marketKey = strings.TrimSpace(filter.Slug)
	}
	if resolvedMarketID == 0 && marketKey != "" {
		var market models.Market
		marketLookup := r.conn.Where("status != 'D'")
		if parsedID, err := strconv.Atoi(marketKey); err == nil && parsedID > 0 {
			marketLookup = marketLookup.Where("id = ? OR slug = ? OR code = ? OR name = ?", parsedID, marketKey, marketKey, marketKey)
		} else {
			marketLookup = marketLookup.Where("slug = ? OR code = ? OR name = ?", marketKey, marketKey, marketKey)
		}
		if err := marketLookup.First(&market).Error; err == nil {
			resolvedMarketID = market.ID
			resolvedMarketSlug = market.Slug
			resolvedMarketName = market.Name
		}
	}
	if resolvedMarketID > 0 {
		query = query.Where("market_id = ?", resolvedMarketID)
	}
	if marketKey != "" {
		if resolvedMarketID == 0 {
			query = query.Where("market_slug = ? OR market_name LIKE ?", marketKey, "%"+marketKey+"%")
		}
	}
	if filter.ZoneID > 0 {
		query = query.Where("zone_id = ?", filter.ZoneID)
	}
	if filter.Search != "" {
		query = query.Where("ident LIKE ? OR name LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	if err := query.Order("ident ASC").Find(&stalls).Error; err != nil {
		return nil, err
	}

	// GIS creates selectable stall codes dynamically, so a booked stall may not
	// have a physical row in the representative seed data. Merge active booking
	// states into the response so the map can still render those stalls as
	// reserved/occupied. No tenant information is exposed here.
	var activeBookings []models.Booking
	bookingQuery := r.conn.Model(&models.Booking{}).
		Where("status = 'A' AND booking_status NOT IN ('waiting_queue', 'rejected', 'not_interested', 'cancelled', 'no_show', 'terminated')")
	if resolvedMarketID > 0 {
		if resolvedMarketName == "" {
			var market models.Market
			if err := r.conn.Where("id = ?", resolvedMarketID).First(&market).Error; err == nil {
				resolvedMarketSlug = market.Slug
				resolvedMarketName = market.Name
			}
		}
		bookingQuery = bookingQuery.Where("market_id = ? OR market_name LIKE ?", resolvedMarketID, "%"+resolvedMarketName+"%")
	} else if marketKey != "" {
		bookingQuery = bookingQuery.Where("market_name LIKE ?", "%"+marketKey+"%")
	}
	if err := bookingQuery.Order("id ASC").Find(&activeBookings).Error; err != nil {
		return nil, err
	}

	stallIndexes := make(map[string]int, len(stalls))
	for i := range stalls {
		stallIndexes[stalls[i].Ident] = i
	}
	for _, booking := range activeBookings {
		ident := strings.TrimSpace(booking.StallIdent)
		if ident == "" {
			continue
		}
		derivedStatus := "reserved"
		if booking.BookingStatus == "lease_applied" || booking.BookingStatus == "completed" {
			derivedStatus = "occupied"
		}

		if index, exists := stallIndexes[ident]; exists {
			// Occupied is stronger than reserved if legacy data contains more than
			// one active workflow for the same stall.
			if stalls[index].StallStatus != "occupied" || derivedStatus == "occupied" {
				stalls[index].StallStatus = derivedStatus
			}
			continue
		}
		if filter.ZoneID > 0 || (filter.Search != "" && !strings.Contains(strings.ToLower(ident), strings.ToLower(filter.Search))) {
			continue
		}

		marketID := resolvedMarketID
		if booking.MarketID != nil {
			marketID = *booking.MarketID
		}
		stalls = append(stalls, models.Stall{
			MarketID:    marketID,
			MarketSlug:  resolvedMarketSlug,
			MarketName:  booking.MarketName,
			ProjectName: booking.ProjectName,
			Ident:       ident,
			Name:        "แผงค้า " + ident,
			StallStatus: derivedStatus,
			Status:      "A",
		})
		stallIndexes[ident] = len(stalls) - 1
	}

	// An active lease is the strongest source of truth for availability. Some
	// GIS stalls do not have a persisted stall row (their geometry lives in the
	// market layout), so merge contracts into the public stall-status response
	// just like bookings. This prevents the booking map from presenting an
	// actively leased stall as available.
	var activeContracts []models.Contract
	contractQuery := r.conn.Model(&models.Contract{}).
		Where("status = 'A' AND contract_state = 'active'")
	if resolvedMarketID > 0 {
		contractQuery = contractQuery.Where("market_id = ?", resolvedMarketID)
	} else if marketKey != "" {
		// A market slug/name has normally been resolved to an ID above. If it
		// cannot be resolved, do not leak statuses from another market.
		contractQuery = contractQuery.Where("1 = 0")
	}
	if err := contractQuery.Order("id ASC").Find(&activeContracts).Error; err != nil {
		return nil, err
	}
	for _, contract := range activeContracts {
		for _, rawIdent := range strings.Split(contract.StallIdent, ",") {
			ident := strings.TrimSpace(rawIdent)
			if ident == "" {
				continue
			}
			if index, exists := stallIndexes[ident]; exists {
				stalls[index].StallStatus = "occupied"
				continue
			}
			if filter.ZoneID > 0 || (filter.Search != "" && !strings.Contains(strings.ToLower(ident), strings.ToLower(filter.Search))) {
				continue
			}

			marketID := resolvedMarketID
			if contract.MarketID != nil {
				marketID = *contract.MarketID
			}
			stalls = append(stalls, models.Stall{
				MarketID:    marketID,
				MarketSlug:  resolvedMarketSlug,
				MarketName:  resolvedMarketName,
				Ident:       ident,
				Name:        "แผงค้า " + ident,
				StallStatus: "occupied",
				Status:      "A",
			})
			stallIndexes[ident] = len(stalls) - 1
		}
	}

	// Merge active waitlist offers ('offered' and not expired) so stalls in the
	// 48-hour waitlist offer window report StallStatus = 'offered'
	var activeOffers []models.Waitlist
	offerQuery := r.conn.Model(&models.Waitlist{}).Where("status = 'offered' AND expires_at > ?", time.Now())
	if resolvedMarketID > 0 {
		offerQuery = offerQuery.Where("market_id = ? OR market_name LIKE ?", resolvedMarketID, "%"+resolvedMarketName+"%")
	} else if marketKey != "" {
		offerQuery = offerQuery.Where("market_name LIKE ?", "%"+marketKey+"%")
	}
	if err := offerQuery.Find(&activeOffers).Error; err == nil {
		for _, offer := range activeOffers {
			ident := strings.TrimSpace(offer.StallIdent)
			if ident == "" {
				continue
			}
			if index, exists := stallIndexes[ident]; exists {
				if stalls[index].StallStatus == "available" || stalls[index].StallStatus == "" {
					stalls[index].StallStatus = "offered"
				}
				continue
			}
			marketID := resolvedMarketID
			if offer.MarketID != nil {
				marketID = *offer.MarketID
			}
			stalls = append(stalls, models.Stall{
				MarketID:    marketID,
				MarketSlug:  resolvedMarketSlug,
				MarketName:  offer.MarketName,
				ProjectName: offer.ProjectName,
				Ident:       ident,
				Name:        "แผงค้า " + ident,
				StallStatus: "offered",
				Status:      "A",
			})
			stallIndexes[ident] = len(stalls) - 1
		}
	}

	if filter.Status == "" {
		return stalls, nil
	}
	filtered := make([]models.Stall, 0, len(stalls))
	for _, stall := range stalls {
		if stall.StallStatus == filter.Status {
			filtered = append(filtered, stall)
		}
	}
	return filtered, nil
}

func (r *repository) FindByID(id int) (*models.Stall, error) {
	var stall models.Stall
	err := r.conn.Where("id = ? AND status = 'A'", id).First(&stall).Error
	if err != nil {
		return nil, err
	}
	return &stall, nil
}

func (r *repository) ByIdent(ident string) (*models.Stall, error) {
	var stall models.Stall
	err := r.conn.Where("ident = ? AND status = 'A'", ident).First(&stall).Error
	if err != nil {
		return nil, err
	}
	return &stall, nil
}

func (r *repository) Create(stall *models.Stall) error {
	now := time.Now()
	stall.CreatedAt = now
	stall.UpdatedAt = now
	stall.Status = "A"
	return r.conn.Transaction(func(tx *gorm.DB) error {
		var existing models.Stall
		if err := tx.Where("ident = ?", stall.Ident).First(&existing).Error; err == nil {
			stall.ID = existing.ID
			stall.CreatedAt = existing.CreatedAt
			if err := tx.Save(stall).Error; err != nil {
				return err
			}
			return tx.Create(&models.StallAuditLog{StallID: stall.ID, Action: "created", Changes: "สร้างแผงเช่า", ActorID: stall.CreatedBy, ActorName: r.actorName(tx, stall.CreatedBy)}).Error
		}
		if err := tx.Create(stall).Error; err != nil {
			return err
		}
		return tx.Create(&models.StallAuditLog{StallID: stall.ID, Action: "created", Changes: "สร้างแผงเช่า", ActorID: stall.CreatedBy, ActorName: r.actorName(tx, stall.CreatedBy)}).Error
	})
}

func (r *repository) actorName(tx *gorm.DB, actorID *int) string {
	if actorID == nil {
		return "ระบบ"
	}
	var user models.User
	if err := tx.Select("full_name", "username").First(&user, *actorID).Error; err != nil {
		return fmt.Sprintf("ผู้ใช้ #%d", *actorID)
	}
	if strings.TrimSpace(user.FullName) != "" {
		return user.FullName
	}
	return user.Username
}

func (r *repository) Update(stall *models.Stall, actorID int, changes string) error {
	return r.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(stall).Error; err != nil {
			return err
		}
		return tx.Create(&models.StallAuditLog{StallID: stall.ID, Action: "updated", Changes: changes, ActorID: &actorID, ActorName: r.actorName(tx, &actorID)}).Error
	})
}

func (r *repository) FindHistory(stallID int) ([]models.StallAuditLog, error) {
	var logs []models.StallAuditLog
	err := r.conn.Where("stall_id = ?", stallID).Order("created_at DESC, id DESC").Find(&logs).Error
	return logs, err
}

func (r *repository) GetActiveBooking(stallIdent string) (*models.Booking, error) {
	var booking models.Booking
	err := r.conn.Where("stall_ident = ? AND queue_number = 0 AND booking_status NOT IN ('waiting_queue', 'rejected', 'not_interested', 'cancelled', 'no_show', 'terminated') AND status = 'A'", stallIdent).
		Order("id ASC").
		First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *repository) GetActiveWaitlistByStall(stallIdent string) ([]models.Waitlist, error) {
	now := time.Now()
	// Expire any stale offers on this stall before returning queue
	_ = r.conn.Model(&models.Waitlist{}).
		Where("stall_ident = ? AND status = 'offered' AND expires_at IS NOT NULL AND expires_at <= ?", strings.TrimSpace(stallIdent), now).
		Updates(map[string]interface{}{"status": "expired", "updated_at": now}).Error

	var items []models.Waitlist
	err := r.conn.Where("stall_ident = ? AND status IN ('offered', 'waiting')", strings.TrimSpace(stallIdent)).
		Order("CASE WHEN status = 'offered' THEN 0 ELSE 1 END ASC, position ASC, id ASC").
		Find(&items).Error
	return items, err
}

// GetActiveWaitlists returns the queue for a market in one query. The list page
// must not call GetActiveWaitlistByStall once per stall because that turns a
// 500-stall market into 500+ database round trips during every map refresh.
func (r *repository) GetActiveWaitlists(filter requests.StallFilter) ([]models.Waitlist, error) {
	now := time.Now()
	expireQuery := r.conn.Model(&models.Waitlist{}).
		Where("status = 'offered' AND expires_at IS NOT NULL AND expires_at <= ?", now)
	query := r.conn.Model(&models.Waitlist{}).
		Where("status IN ('offered', 'waiting')")
	if filter.MarketID > 0 {
		expireQuery = expireQuery.Where("market_id = ?", filter.MarketID)
		query = query.Where("market_id = ?", filter.MarketID)
	}
	_ = expireQuery.Updates(map[string]interface{}{"status": "expired", "updated_at": now}).Error

	var items []models.Waitlist
	err := query.
		Order("stall_ident ASC, CASE WHEN status = 'offered' THEN 0 ELSE 1 END ASC, position ASC, id ASC").
		Find(&items).Error
	return items, err
}

func (r *repository) CountQueue(stallIdent string) (int64, error) {
	var total int64
	err := r.conn.Model(&models.Waitlist{}).
		Where("stall_ident = ? AND status IN ('waiting', 'offered')", strings.TrimSpace(stallIdent)).
		Count(&total).Error
	return total, err
}
