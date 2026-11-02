package databases

import (
	"log"
	"strings"

	"emarket-api/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitSQLite(dbPath string) (*gorm.DB, error) {
	var err error
	// WAL lets reads continue while a write is in progress, and busy_timeout makes
	// a blocked writer wait instead of failing with "database is locked".
	dsn := dbPath
	if dbPath != "" && !strings.Contains(dbPath, ":memory:") && !strings.Contains(dbPath, "_pragma=") {
		sep := "?"
		if strings.Contains(dbPath, "?") {
			sep = "&"
		}
		dsn = dbPath + sep + "_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)"
	}
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	log.Println("SQLite connected successfully:", dbPath)
	return DB, nil
}

func AutoMigrate(db *gorm.DB) error {
	// Normalize legacy 'submitted' states to 'pending_officer'
	_ = db.Model(&models.Application{}).Where("application_state = ? OR application_state = ''", "submitted").Update("application_state", "pending_officer").Error

	if err := db.AutoMigrate(
		&models.User{},
		&models.MarketStaffAssignment{},
		&models.Market{},
		&models.Zone{},
		&models.Stall{},
		&models.StallAuditLog{},
		&models.MarketGISConfig{},
		&models.MarketUtilitySetting{},
		&models.MarketTradeCondition{},
		&models.MarketContractTemplate{},
		&models.Booking{},
		&models.BookingStall{},
		&models.Waitlist{},
		&models.Notification{},
		&models.Application{},
		&models.ApplicationType{},
		&models.Debt{},
		&models.MaintenanceRecord{},
		&models.Contract{},
		&models.Item{},
	); err != nil {
		return err
	}

	var jjTemplateCount int64
	db.Model(&models.MarketContractTemplate{}).Where("market_id = ? AND code = ?", 1, "jj_semi_permanent_lease").Count(&jjTemplateCount)
	if jjTemplateCount == 0 {
		_ = db.Create(&models.MarketContractTemplate{
			MarketID:     1,
			Code:         "jj_semi_permanent_lease",
			Title:        "สัญญาให้เช่าแผงค้าภายในบริเวณตลาดนัดจตุจักร (แผงค้ากึ่งถาวร)",
			ContractKind: "standard",
			Description:  "บันทึกข้อตกลงสัญญาให้เช่าแผงค้ากึ่งถาวรภายในบริเวณตลาดนัดจตุจักร สำนักงานตลาดกรุงเทพมหานคร กำหนดระยะเวลา 1 ปี",
			IsActive:     true,
			SortOrder:    1,
		}).Error
	}

	return nil
}
