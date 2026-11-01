package repositories

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"emarket-api/models"
	"emarket-api/src/market/domains"
	"gorm.io/gorm"
)

type repository struct {
	conn *gorm.DB
}

func NewRepository(db *gorm.DB) domains.Repository {
	return &repository{conn: db}
}

func (r *repository) syncMarketStats(m *models.Market) {
	var occupiedCount int64
	r.conn.Model(&models.Stall{}).
		Where("status = 'A' AND (market_id = ? OR market_name LIKE ?) AND stall_status = 'occupied'", m.ID, "%"+m.Name+"%").
		Distinct("ident").
		Count(&occupiedCount)

	var activeContractCount int64
	r.conn.Model(&models.Contract{}).
		Where("status = 'A' AND contract_state = 'active' AND stall_ident IN (?)",
			r.conn.Model(&models.Stall{}).Select("ident").Where("market_id = ? OR market_name LIKE ?", m.ID, "%"+m.Name+"%"),
		).
		Distinct("stall_ident").
		Count(&activeContractCount)
	if activeContractCount > occupiedCount {
		occupiedCount = activeContractCount
	}

	var reservedCount int64
	r.conn.Model(&models.Stall{}).
		Where("status = 'A' AND (market_id = ? OR market_name LIKE ?) AND stall_status = 'reserved'", m.ID, "%"+m.Name+"%").
		Distinct("ident").
		Count(&reservedCount)

	var activeBookingReserved int64
	r.conn.Model(&models.Booking{}).
		Where("status = 'A' AND (market_id = ? OR market_name LIKE ?) AND booking_status NOT IN ('waiting_queue', 'rejected', 'not_interested', 'cancelled', 'no_show', 'lease_applied', 'completed', 'terminated', 'closed')", m.ID, "%"+m.Name+"%").
		Distinct("stall_ident").
		Count(&activeBookingReserved)
	if activeBookingReserved > reservedCount {
		reservedCount = activeBookingReserved
	}

	m.OccupiedStalls = int(occupiedCount)
	m.ReservedStalls = int(reservedCount)
	avail := m.TotalStalls - m.OccupiedStalls - m.ReservedStalls
	if avail < 0 {
		avail = 0
	}
	m.AvailableStalls = avail
	if m.TotalStalls > 0 && m.AvailableStalls == 0 {
		m.Status = "full"
	}
}

type marketStallStatusRow struct {
	MarketID    int
	MarketName  string
	Ident       string
	StallStatus string
}

type marketContractStatusRow struct {
	MarketID   *int
	StallIdent string
}

type marketBookingStatusRow struct {
	MarketID   *int
	MarketName string
	StallIdent string
}

type marketStatusSets struct {
	occupied map[string]struct{}
	reserved map[string]struct{}
}

// syncMarketsStats calculates the list-page counters in three bounded queries.
// The previous implementation ran four aggregate queries for every market,
// which made the public tenant market page wait on dozens of full-table scans.
func (r *repository) syncMarketsStats(markets []models.Market) error {
	if len(markets) == 0 {
		return nil
	}

	statsByMarketID := make(map[int]*marketStatusSets, len(markets))
	marketIDByName := make(map[string]int, len(markets))
	for i := range markets {
		statsByMarketID[markets[i].ID] = &marketStatusSets{
			occupied: make(map[string]struct{}),
			reserved: make(map[string]struct{}),
		}
		marketIDByName[strings.ToLower(strings.TrimSpace(markets[i].Name))] = markets[i].ID
	}

	resolveMarketID := func(marketID int, marketName string) int {
		if _, ok := statsByMarketID[marketID]; ok {
			return marketID
		}
		name := strings.ToLower(strings.TrimSpace(marketName))
		if id, ok := marketIDByName[name]; ok {
			return id
		}
		// Preserve support for legacy rows whose market_name contains a prefix
		// or suffix around the canonical market name.
		for canonicalName, id := range marketIDByName {
			if canonicalName != "" && strings.Contains(name, canonicalName) {
				return id
			}
		}
		return 0
	}

	var stalls []marketStallStatusRow
	if err := r.conn.Model(&models.Stall{}).
		Select("market_id, market_name, ident, stall_status").
		Where("status = 'A'").
		Find(&stalls).Error; err != nil {
		return err
	}

	stallMarketByIdent := make(map[string]int, len(stalls))
	for _, stall := range stalls {
		marketID := resolveMarketID(stall.MarketID, stall.MarketName)
		stats := statsByMarketID[marketID]
		if stats == nil || stall.Ident == "" {
			continue
		}
		stallMarketByIdent[stall.Ident] = marketID
		switch strings.ToLower(strings.TrimSpace(stall.StallStatus)) {
		case "occupied":
			stats.occupied[stall.Ident] = struct{}{}
		case "reserved":
			stats.reserved[stall.Ident] = struct{}{}
		}
	}

	var contracts []marketContractStatusRow
	if err := r.conn.Model(&models.Contract{}).
		Select("market_id, stall_ident").
		Where("status = 'A' AND contract_state = 'active'").
		Find(&contracts).Error; err != nil {
		return err
	}
	for _, contract := range contracts {
		marketID := 0
		if contract.MarketID != nil {
			marketID = *contract.MarketID
		}
		if _, ok := statsByMarketID[marketID]; !ok {
			marketID = stallMarketByIdent[contract.StallIdent]
		}
		if stats := statsByMarketID[marketID]; stats != nil && contract.StallIdent != "" {
			stats.occupied[contract.StallIdent] = struct{}{}
		}
	}

	var bookings []marketBookingStatusRow
	if err := r.conn.Model(&models.Booking{}).
		Select("market_id, market_name, stall_ident, booking_status").
		Where("status = 'A' AND booking_status NOT IN ?", []string{
			"waiting_queue", "rejected", "not_interested", "cancelled", "no_show",
			"lease_applied", "completed", "terminated", "closed",
		}).
		Find(&bookings).Error; err != nil {
		return err
	}
	for _, booking := range bookings {
		marketID := 0
		if booking.MarketID != nil {
			marketID = *booking.MarketID
		}
		marketID = resolveMarketID(marketID, booking.MarketName)
		if stats := statsByMarketID[marketID]; stats != nil && booking.StallIdent != "" {
			stats.reserved[booking.StallIdent] = struct{}{}
		}
	}

	for i := range markets {
		stats := statsByMarketID[markets[i].ID]
		if stats == nil {
			continue
		}
		// A contracted stall is occupied even if an older booking/status row was
		// not cleaned up yet; never count the same stall in both counters.
		for ident := range stats.occupied {
			delete(stats.reserved, ident)
		}
		markets[i].OccupiedStalls = len(stats.occupied)
		markets[i].ReservedStalls = len(stats.reserved)
		available := markets[i].TotalStalls - markets[i].OccupiedStalls - markets[i].ReservedStalls
		if available < 0 {
			available = 0
		}
		markets[i].AvailableStalls = available
		if markets[i].TotalStalls > 0 && available == 0 {
			markets[i].Status = "full"
		}
	}

	return nil
}

func (r *repository) FindAll(search string) ([]models.Market, error) {
	var markets []models.Market
	query := r.conn.Model(&models.Market{}).Where("status != 'D'")
	if search != "" {
		s := "%" + search + "%"
		query = query.Where("name LIKE ? OR name_en LIKE ? OR district LIKE ? OR code LIKE ?", s, s, s, s)
	}
	err := query.Order("id ASC").Find(&markets).Error
	if err != nil {
		return nil, err
	}
	if err := r.syncMarketsStats(markets); err != nil {
		return nil, err
	}
	return markets, nil
}

func (r *repository) FindByStaffUserID(search string, userID int) ([]models.Market, error) {
	var markets []models.Market
	query := r.conn.Model(&models.Market{}).
		Joins("JOIN market_staff_assignment msa ON msa.market_id = market.id").
		Where("msa.user_id = ? AND market.status != 'D'", userID)
	if search != "" {
		s := "%" + search + "%"
		query = query.Where("market.name LIKE ? OR market.name_en LIKE ? OR market.district LIKE ? OR market.code LIKE ?", s, s, s, s)
	}
	if err := query.Order("market.id ASC").Find(&markets).Error; err != nil {
		return nil, err
	}
	if err := r.syncMarketsStats(markets); err != nil {
		return nil, err
	}
	return markets, nil
}

func (r *repository) FindByIDOrSlug(idOrSlug string) (*models.Market, error) {
	var market models.Market
	query := r.conn.Where("status != 'D'")
	if id, err := strconv.Atoi(idOrSlug); err == nil {
		query = query.Where("id = ? OR slug = ?", id, idOrSlug)
	} else {
		query = query.Where("slug = ?", idOrSlug)
	}

	err := query.First(&market).Error
	if err != nil {
		return nil, err
	}
	r.syncMarketStats(&market)
	return &market, nil
}

func (r *repository) FindZonesByMarketID(marketID int) ([]models.Zone, error) {
	var zones []models.Zone
	query := r.conn.Where("status = 'A'")
	if marketID > 0 {
		query = query.Where("market_id = ?", marketID)
	}
	err := query.Order("id ASC").Find(&zones).Error
	return zones, err
}

func (r *repository) FindGISConfig(marketID int) (*models.MarketGISConfig, error) {
	var config models.MarketGISConfig
	if err := r.conn.Where("market_id = ?", marketID).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// FindGISConfigMeta reads only the version columns so callers can decide whether a
// cached copy of the (multi-MB) projects JSON is still current.
func (r *repository) FindGISConfigMeta(marketID int) (updatedAt time.Time, published bool, err error) {
	var config models.MarketGISConfig
	if err = r.conn.Select("id", "published", "updated_at").Where("market_id = ?", marketID).First(&config).Error; err != nil {
		return time.Time{}, false, err
	}
	return config.UpdatedAt, config.Published, nil
}

func (r *repository) SaveGISConfig(config *models.MarketGISConfig) error {
	var existing models.MarketGISConfig
	if err := r.conn.Where("market_id = ?", config.MarketID).First(&existing).Error; err == nil {
		type projectItem struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Stalls []struct {
				Code string `json:"code"`
			} `json:"stalls"`
		}
		var oldProjects, newProjects []projectItem
		_ = json.Unmarshal([]byte(existing.ProjectsJSON), &oldProjects)
		_ = json.Unmarshal([]byte(config.ProjectsJSON), &newProjects)

		activeProjectIDs := make(map[string]bool, len(newProjects))
		activeProjectNames := make(map[string]bool, len(newProjects))
		for _, p := range newProjects {
			activeProjectIDs[p.ID] = true
			activeProjectNames[p.Name] = true
		}

		var removedProjectNames []string
		var removedStallCodes []string
		for _, oldP := range oldProjects {
			if !activeProjectIDs[oldP.ID] && !activeProjectNames[oldP.Name] {
				if oldP.Name != "" {
					removedProjectNames = append(removedProjectNames, oldP.Name)
				}
				for _, s := range oldP.Stalls {
					if s.Code != "" {
						removedStallCodes = append(removedStallCodes, s.Code)
					}
				}
			}
		}

		if len(removedProjectNames) > 0 {
			r.conn.Model(&models.Stall{}).
				Where("market_id = ? AND project_name IN ?", config.MarketID, removedProjectNames).
				Update("status", "D")
		}
		if len(removedStallCodes) > 0 {
			r.conn.Model(&models.Stall{}).
				Where("market_id = ? AND ident IN ?", config.MarketID, removedStallCodes).
				Update("status", "D")
		}
	}

	return r.conn.Where(models.MarketGISConfig{MarketID: config.MarketID}).
		Assign(map[string]interface{}{"projects_json": config.ProjectsJSON, "facilities_json": config.FacilitiesJSON, "published": config.Published, "updated_by": config.UpdatedBy, "updated_at": time.Now()}).
		FirstOrCreate(config).Error
}

func (r *repository) FindUtilitySettings(marketID int) (*models.MarketUtilitySetting, error) {
	var setting models.MarketUtilitySetting
	if err := r.conn.Where("market_id = ?", marketID).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// SaveUtilitySettings แทนที่เอกสารตั้งค่าของตลาด (หนึ่งแถวต่อตลาด)
func (r *repository) SaveUtilitySettings(setting *models.MarketUtilitySetting) error {
	var existing models.MarketUtilitySetting
	if err := r.conn.Where("market_id = ?", setting.MarketID).First(&existing).Error; err == nil {
		return r.conn.Model(&existing).Updates(map[string]interface{}{
			"settings_json": setting.SettingsJSON,
			"updated_by":    setting.UpdatedBy,
			"updated_at":    time.Now(),
		}).Error
	}
	return r.conn.Create(setting).Error
}
