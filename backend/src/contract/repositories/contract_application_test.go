package repositories

import (
	"errors"
	"testing"

	"emarket-api/models"
	"emarket-api/src/contract/domains"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestCreateContractRequiresApprovedMatchingApplication(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Application{}, &models.Contract{}, &models.Stall{}, &models.Booking{}, &models.Waitlist{}, &models.Notification{}); err != nil {
		t.Fatal(err)
	}
	app := models.Application{ApplicationNo: "APP-CONTRACT-TEST", StallIdent: "A-1", ApplicantName: "Tenant", ApplicantPhone: "1", ShopName: "Shop", ApplicationState: "pending_bureau_director", Status: "A"}
	if err := db.Create(&app).Error; err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)
	contract := func(ident string) *models.Contract {
		return &models.Contract{ContractNumber: "CTR-TEST", ApplicationID: &app.ID, StallIdent: ident, TenantName: "Tenant"}
	}
	if err := repo.Create(contract("A-1")); !errors.Is(err, domains.ErrApplicationNotApproved) {
		t.Fatalf("pending application: got %v", err)
	}
	if err := db.Model(&app).Update("application_state", "approved").Error; err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(contract("B-2")); !errors.Is(err, domains.ErrApplicationNotApproved) {
		t.Fatalf("different stall: got %v", err)
	}
	if err := repo.Create(contract("A-1")); err != nil {
		t.Fatalf("approved application: %v", err)
	}
	if err := db.First(&app, app.ID).Error; err != nil {
		t.Fatal(err)
	}
	if app.ApplicationState != "completed" {
		t.Fatalf("application state = %q", app.ApplicationState)
	}
	if err := repo.Create(contract("A-1")); !errors.Is(err, domains.ErrApplicationNotApproved) {
		t.Fatalf("duplicate contract: got %v", err)
	}
}
