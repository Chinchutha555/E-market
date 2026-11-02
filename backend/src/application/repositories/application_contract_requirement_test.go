package repositories

import (
	"testing"

	"emarket-api/models"
	"emarket-api/requests"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestApplicationReadsBookingContractRequirement(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:application-contract-requirement?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Booking{}, &models.Application{}); err != nil {
		t.Fatal(err)
	}
	noContract := false
	booking := models.Booking{BookingNo: "REQ-BOOKING", StallIdent: "S-1", ApplicantName: "Tenant", ApplicantPhone: "123", VisitDate: "2026-10-01", VisitTime: "10:00", RequiresContract: &noContract}
	if err := db.Create(&booking).Error; err != nil {
		t.Fatal(err)
	}
	app := models.Application{ApplicationNo: "REQ-APPLICATION", ApplicationType: "stall_lease", BookingID: &booking.ID, StallIdent: "S-1", ApplicantName: "Tenant", ApplicantPhone: "123", ShopName: "Shop", Status: "A"}
	if err := db.Create(&app).Error; err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	found, err := repo.FindByID(app.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.RequiresContract == nil || *found.RequiresContract {
		t.Fatal("detail did not preserve no-contract requirement")
	}
	list, err := repo.FindAll(requests.ApplicationFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].RequiresContract == nil || *list[0].RequiresContract {
		t.Fatal("list did not preserve no-contract requirement")
	}
}
