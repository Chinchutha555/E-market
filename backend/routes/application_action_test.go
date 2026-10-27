package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"emarket-api/config"
	"emarket-api/databases"
	"emarket-api/helpers"
	"emarket-api/models"
	"emarket-api/requests"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestApplicationRejectWithComment(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open("file:app-reject-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite memory db: %v", err)
	}
	databases.DB = db
	if err := databases.AutoMigrate(db); err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}

	market := models.Market{ID: 1, Slug: "mkt-test", Code: "MKT-01", Name: "ตลาดทดสอบ", Status: "open"}
	if err := db.Create(&market).Error; err != nil {
		t.Fatalf("failed to create market: %v", err)
	}

	stall := models.Stall{ID: 1, Ident: "A-101", MarketID: market.ID, Status: "A"}
	if err := db.Create(&stall).Error; err != nil {
		t.Fatalf("failed to create stall: %v", err)
	}

	staff := models.User{
		ID:       55,
		Username: "officer-test",
		Password: "x",
		FullName: "เจ้าหน้าที่ทดสอบ",
		Role:     "staff",
		Status:   "A",
	}
	if err := db.Create(&staff).Error; err != nil {
		t.Fatalf("failed to create staff: %v", err)
	}

	assignment := models.MarketStaffAssignment{UserID: staff.ID, MarketID: market.ID}
	if err := db.Create(&assignment).Error; err != nil {
		t.Fatalf("failed to create assignment: %v", err)
	}

	app := models.Application{
		ApplicationNo:    "APP-TEST-REJECT-001",
		StallIdent:       "A-101",
		MarketID:         &market.ID,
		ApplicantName:    "สมชาย ทดสอบ",
		ApplicantPhone:   "0899999999",
		ShopName:         "ร้านค้าทดสอบ",
		ApplicationState: "pending_officer",
		Status:           "A",
	}
	if err := db.Create(&app).Error; err != nil {
		t.Fatalf("failed to create test application: %v", err)
	}

	token, err := helpers.GenerateToken(staff.ID, staff.Username, staff.Role)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	router := SetupRouter()

	rejectReason := "เอกสารไม่ครบถ้วน บัตรประชาชนหมดอายุ"
	payload := requests.ApplicationActionRequest{
		Action:  "reject",
		Comment: rejectReason,
		Reason:  rejectReason,
	}
	payloadBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/applications/"+strconv.Itoa(app.ID)+"/action", bytes.NewReader(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d: %s", res.Code, res.Body.String())
	}

	var jsonRes struct {
		Status bool `json:"status"`
		Data   struct {
			ID             int    `json:"id"`
			Status         string `json:"status"`
			OfficerComment string `json:"officer_comment"`
			ReviewNote     string `json:"review_note"`
		} `json:"data"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &jsonRes); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}

	if jsonRes.Data.Status != "rejected" {
		t.Errorf("expected status 'rejected', got '%s'", jsonRes.Data.Status)
	}
	if jsonRes.Data.OfficerComment != rejectReason {
		t.Errorf("expected officer_comment '%s', got '%s'", rejectReason, jsonRes.Data.OfficerComment)
	}
	if jsonRes.Data.ReviewNote != rejectReason {
		t.Errorf("expected review_note '%s', got '%s'", rejectReason, jsonRes.Data.ReviewNote)
	}

	// Verify persistence in SQLite DB
	var dbApp models.Application
	if err := db.First(&dbApp, app.ID).Error; err != nil {
		t.Fatalf("failed to find updated application in DB: %v", err)
	}
	if dbApp.ApplicationState != "rejected" {
		t.Errorf("expected DB ApplicationState 'rejected', got '%s'", dbApp.ApplicationState)
	}
	if dbApp.OfficerComment != rejectReason {
		t.Errorf("expected DB OfficerComment '%s', got '%s'", rejectReason, dbApp.OfficerComment)
	}
	if dbApp.ReviewNote != rejectReason {
		t.Errorf("expected DB ReviewNote '%s', got '%s'", rejectReason, dbApp.ReviewNote)
	}
}
