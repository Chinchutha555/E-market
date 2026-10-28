package databases

import (
	"log"

	"emarket-api/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitSQLite(dbPath string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	log.Println("SQLite connected successfully:", dbPath)
	return DB, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.MarketStaffAssignment{},
		&models.Market{},
		&models.Zone{},
		&models.Stall{},
		&models.StallAuditLog{},
		&models.MarketGISConfig{},
		&models.Booking{},
		&models.BookingStall{},
		&models.Waitlist{},
		&models.Notification{},
		&models.Application{},
		&models.ApplicationType{},
		&models.Debt{},
		&models.Contract{},
		&models.Item{},
	)
}
