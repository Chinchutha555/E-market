package repositories

import (
	"testing"

	"emarket-api/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestCreateApplicationClosesActiveWaitlistAndNotifiesTenants(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(
		&models.Application{},
		&models.Booking{},
		&models.Stall{},
		&models.Waitlist{},
		&models.Notification{},
	); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	tenantID := 101
	holder := models.Booking{
		BookingNo:          "BK-TWR-PRIMARY",
		StallIdent:         "TWR-03-026",
		ApplicantName:      "ผู้เช่าหลัก",
		ApplicantPhone:     "0800000000",
		VisitDate:          "2026-10-26",
		VisitTime:          "09:00",
		BookingStatus:      "interested",
		AppointmentGroupID: "BKG-TWR-PRIMARY",
		Status:             "A",
		CreatedBy:          &tenantID,
	}
	if err := db.Create(&holder).Error; err != nil {
		t.Fatalf("create holder: %v", err)
	}

	waiters := []models.Waitlist{
		{StallIdent: "TWR-03-026", TenantID: 201, TenantName: "ผู้รอหนึ่ง", Position: 1, Status: "waiting"},
		{StallIdent: "TWR-03-026", TenantID: 202, TenantName: "ผู้รอสอง", Position: 2, Status: "waiting"},
	}
	if err := db.Create(&waiters).Error; err != nil {
		t.Fatalf("create waitlist: %v", err)
	}

	app := &models.Application{
		ApplicationNo:    "APP-TWR-026",
		StallIdent:       "TWR-03-026",
		ApplicantName:    "ผู้เช่าหลัก",
		ApplicantPhone:   "0800000000",
		ShopName:         "ร้านทดสอบ",
		ApplicationState: "submitted",
		CreatedBy:        &tenantID,
	}
	if err := NewRepository(db).Create(app); err != nil {
		t.Fatalf("create application: %v", err)
	}

	var updatedHolder models.Booking
	if err := db.First(&updatedHolder, holder.ID).Error; err != nil {
		t.Fatalf("reload holder: %v", err)
	}
	if updatedHolder.BookingStatus != "lease_applied" {
		t.Fatalf("holder status = %q, want lease_applied", updatedHolder.BookingStatus)
	}

	var activeCount int64
	if err := db.Model(&models.Waitlist{}).
		Where("stall_ident = ? AND status IN ('waiting', 'offered')", "TWR-03-026").
		Count(&activeCount).Error; err != nil {
		t.Fatalf("count active waitlist: %v", err)
	}
	if activeCount != 0 {
		t.Fatalf("active waitlist count = %d, want 0", activeCount)
	}

	var cancelledCount int64
	if err := db.Model(&models.Waitlist{}).
		Where("stall_ident = ? AND status = ?", "TWR-03-026", "cancelled").
		Count(&cancelledCount).Error; err != nil {
		t.Fatalf("count cancelled waitlist: %v", err)
	}
	if cancelledCount != 2 {
		t.Fatalf("cancelled waitlist count = %d, want 2", cancelledCount)
	}

	var notificationCount int64
	if err := db.Model(&models.Notification{}).
		Where("stall_number = ? AND status_code = ?", "TWR-03-026", "waitlist_cancelled").
		Count(&notificationCount).Error; err != nil {
		t.Fatalf("count waitlist notifications: %v", err)
	}
	if notificationCount != 2 {
		t.Fatalf("notification count = %d, want 2", notificationCount)
	}
}
