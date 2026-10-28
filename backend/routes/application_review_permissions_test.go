package routes

import (
	"bytes"
	"emarket-api/config"
	"emarket-api/databases"
	"emarket-api/helpers"
	"emarket-api/models"
	"encoding/json"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestApplicationReviewerPermissions(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	databases.DB = db
	if err = databases.AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	market := models.Market{ID: 1, Slug: "review-test", Code: "RT", Name: "Test", Status: "open"}
	if err = db.Create(&market).Error; err != nil {
		t.Fatal(err)
	}
	users := []models.User{{Username: "officer", Password: "x", Role: "staff", Status: "A"}, {Username: "market-director", Password: "x", Role: "market_director", Status: "A"}, {Username: "bureau-director", Password: "x", Role: "bureau_director", Status: "A"}, {Username: "unassigned-director", Password: "x", Role: "market_director", Status: "A"}}
	if err = db.Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	for _, u := range users[:2] {
		if err = db.Create(&models.MarketStaffAssignment{UserID: u.ID, MarketID: market.ID}).Error; err != nil {
			t.Fatal(err)
		}
	}
	router := SetupRouter()
	app := models.Application{ApplicationNo: "PERMISSION-TEST", ApplicationType: "stall_lease", MarketID: &market.ID, StallIdent: "TEST-P", ApplicantName: "Test", ApplicantPhone: "1", ShopName: "Test", ApplicationState: "pending_officer", Status: "A"}
	if err = db.Create(&app).Error; err != nil {
		t.Fatal(err)
	}
	call := func(user models.User, action string) int {
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		payload, _ := json.Marshal(map[string]string{"action": action, "comment": "checked"})
		req := httptest.NewRequest("POST", "/api/v1/applications/"+strconv.Itoa(app.ID)+"/action", bytes.NewReader(payload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec.Code
	}
	if code := call(users[1], "approve"); code != 403 {
		t.Fatalf("director before officer: %d", code)
	}
	if code := call(users[0], "approve"); code != 200 {
		t.Fatalf("officer: %d", code)
	}
	if code := call(users[0], "approve"); code != 403 {
		t.Fatalf("officer at director step: %d", code)
	}
	if code := call(users[3], "approve"); code != 403 {
		t.Fatalf("unassigned director: %d", code)
	}
	if code := call(users[1], "approve"); code != 200 {
		t.Fatalf("market director: %d", code)
	}
	if code := call(users[1], "approve"); code != 403 {
		t.Fatalf("market director at bureau step: %d", code)
	}
	if code := call(users[2], "approve"); code != 200 {
		t.Fatalf("bureau director: %d", code)
	}
}
