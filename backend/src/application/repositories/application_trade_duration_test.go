package repositories

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"emarket-api/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestApplicationDurationUsesBookingSnapshot(t *testing.T) {
	for _, tc := range []struct {
		unit           string
		amount, months int
	}{{"days", 7, 0}, {"months", 3, 3}, {"years", 2, 24}} {
		t.Run(tc.unit, func(t *testing.T) {
			db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:duration-%s-%d?mode=memory&cache=shared", tc.unit, time.Now().UnixNano())), &gorm.Config{})
			if err != nil {
				t.Fatal(err)
			}
			if err = db.AutoMigrate(&models.Application{}, &models.Booking{}, &models.MarketTradeCondition{}, &models.Stall{}, &models.Waitlist{}, &models.Notification{}); err != nil {
				t.Fatal(err)
			}
			user := 123
			booking := models.Booking{BookingNo: "TERM-BK", StallIdent: "S-1", ApplicantName: "Tenant", ApplicantPhone: "123", VisitDate: "2026-10-01", VisitTime: "10:00", BookingStatus: "interested", Status: "A", RentalType: "custom", TradeDuration: tc.amount, TradeDurationUnit: tc.unit, CreatedBy: &user}
			if err = db.Create(&booking).Error; err != nil {
				t.Fatal(err)
			}
			app := models.Application{ApplicationNo: "TERM-APP", ApplicationType: "stall_lease", BookingID: &booking.ID, StallIdent: "S-1", ApplicantName: "Tenant", ApplicantPhone: "123", ShopName: "Shop", LeaseTermMonths: 99, ProposedRentType: "wrong", Status: "A"}
			if err = NewRepository(db).Create(&app); err != nil {
				t.Fatal(err)
			}
			var saved models.Application
			if err = db.First(&saved, app.ID).Error; err != nil {
				t.Fatal(err)
			}
			if saved.TradeDuration != tc.amount || saved.TradeDurationUnit != tc.unit || saved.LeaseTermMonths != tc.months || saved.ProposedRentType != "custom" {
				t.Fatalf("wrong duration: %#v", saved)
			}
		})
	}
}

func TestApplicationSelectsOnlyBookingDurationOptions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:multiple-duration-%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&models.Application{}, &models.Booking{}, &models.MarketTradeCondition{}, &models.Stall{}, &models.Waitlist{}, &models.Notification{}); err != nil {
		t.Fatal(err)
	}
	options, _ := json.Marshal([]models.TradeDurationOption{
		{ID: "one", Duration: 1, Unit: "months", StartDate: "2026-11-02", EndDate: "2026-12-01"},
		{ID: "three", Duration: 3, Unit: "months", StartDate: "2026-10-01", EndDate: "2026-12-31"},
	})
	creator := 123
	booking := models.Booking{BookingNo: "MULTI-BK", StallIdent: "S-2", ApplicantName: "Tenant", ApplicantPhone: "123", VisitDate: "2026-09-01", VisitTime: "10:00", BookingStatus: "interested", Status: "A", RentalType: "trial", TradeOptionsJSON: string(options), CreatedBy: &creator}
	if err = db.Create(&booking).Error; err != nil {
		t.Fatal(err)
	}
	makeApp := func(id string) *models.Application {
		return &models.Application{ApplicationNo: "MULTI-APP", ApplicationType: "stall_lease", BookingID: &booking.ID, StallIdent: "S-2", ApplicantName: "Tenant", ApplicantPhone: "123", ShopName: "Shop", RequestedStartDate: "2030-01-01", LeaseTermMonths: 99, TradeOptionID: id, Status: "A"}
	}
	repo := NewRepository(db)
	if err = repo.Create(makeApp("")); err == nil {
		t.Fatal("accepted missing choice")
	}
	if err = repo.Create(makeApp("forged")); err == nil {
		t.Fatal("accepted forged choice")
	}
	app := makeApp("three")
	if err = repo.Create(app); err != nil {
		t.Fatal(err)
	}
	var saved models.Application
	if err = db.First(&saved, app.ID).Error; err != nil {
		t.Fatal(err)
	}
	if saved.TradeOptionID != "three" || saved.TradeDuration != 3 || saved.TradeDurationUnit != "months" || saved.LeaseTermMonths != 3 || saved.RequestedStartDate != "" || saved.TradeEndDate != "" {
		t.Fatalf("wrong option: %#v", saved)
	}
}

func TestLegacyBookingUsesCurrentOptions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:legacy-options-%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&models.Application{}, &models.Booking{}, &models.MarketTradeCondition{}, &models.Stall{}, &models.Waitlist{}, &models.Notification{}); err != nil {
		t.Fatal(err)
	}
	options, _ := json.Marshal([]models.TradeDurationOption{{ID: "one", Duration: 1, Unit: "months"}, {ID: "three", Duration: 3, Unit: "months"}})
	market := 1
	setting := models.MarketTradeCondition{MarketID: market, TradeType: "trial", Title: "Trial", MaxDuration: 1, MinDuration: 1, Unit: "months", OptionsJSON: string(options)}
	if err = db.Create(&setting).Error; err != nil {
		t.Fatal(err)
	}
	booking := models.Booking{BookingNo: "LEGACY-BK", StallIdent: "S-3", ApplicantName: "Tenant", ApplicantPhone: "123", VisitDate: "2026-10-01", VisitTime: "10:00", BookingStatus: "interested", Status: "A", RentalType: "trial", MarketID: &market, TradeDuration: 1, TradeDurationUnit: "months"}
	if err = db.Create(&booking).Error; err != nil {
		t.Fatal(err)
	}
	app := models.Application{ApplicationNo: "LEGACY-APP", ApplicationType: "stall_lease", BookingID: &booking.ID, StallIdent: "S-3", ApplicantName: "Tenant", ApplicantPhone: "123", ShopName: "Shop", TradeOptionID: "three", Status: "A"}
	if err = NewRepository(db).Create(&app); err != nil {
		t.Fatal(err)
	}
	if app.TradeDuration != 3 || app.TradeOptionID != "three" {
		t.Fatalf("wrong selected option: %#v", app)
	}
}
