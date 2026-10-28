package routes

import (
	"bytes"
	"net/http/httptest"
	"strconv"
	"testing"

	"emarket-api/config"
	"emarket-api/databases"
	"emarket-api/helpers"
	"emarket-api/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestTenantResubmitReturnedApplication(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	databases.DB = db
	if err := databases.AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	users := []models.User{{Username: "resubmit-owner", Password: "x", Role: "tenant", Status: "A"}, {Username: "resubmit-other", Password: "x", Role: "tenant", Status: "A"}}
	if err := db.Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	app := models.Application{ApplicationNo: "RESUBMIT-TEST", ApplicationType: "stall_lease", StallIdent: "A-1", ApplicantName: "Old", ApplicantPhone: "1", ShopName: "Shop", ApplicationState: "revision_required", Status: "A", CreatedBy: &users[0].ID}
	if err := db.Create(&app).Error; err != nil {
		t.Fatal(err)
	}
	router := SetupRouter()
	call := func(user models.User) int {
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		req := httptest.NewRequest("POST", "/api/v1/applications/"+strconv.Itoa(app.ID)+"/resubmit", bytes.NewBufferString(`{"applicant_name":"Updated","applicant_phone":"123","shop_name":"Shop","details":"Fixed"}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec.Code
	}
	if code := call(users[1]); code != 403 {
		t.Fatalf("other tenant: %d", code)
	}
	if code := call(users[0]); code != 200 {
		t.Fatalf("owner: %d", code)
	}
	if err := db.First(&app, app.ID).Error; err != nil {
		t.Fatal(err)
	}
	if app.ApplicationState != "pending_officer" || app.ApplicantName != "Updated" || app.Details != "Fixed" {
		t.Fatalf("unexpected result: %+v", app)
	}
	if code := call(users[0]); code != 409 {
		t.Fatalf("repeat: %d", code)
	}
}
