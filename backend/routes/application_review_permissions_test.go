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
	admin := models.User{Username: "application-admin", Password: "x", Role: "admin", Status: "A"}
	if err := db.Create(&admin).Error; err != nil {
		t.Fatal(err)
	}
	router := SetupRouter()
	app := models.Application{ApplicationNo: "PERMISSION-TEST", ApplicationType: "stall_lease", MarketID: &market.ID, StallIdent: "TEST-P", ApplicantName: "Test", ApplicantPhone: "1", ShopName: "Test", ApplicationState: "pending_officer", Status: "A"}
	if err = db.Create(&app).Error; err != nil {
		t.Fatal(err)
	}
	otp := ""
	call := func(user models.User, action string) int {
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		payload, _ := json.Marshal(map[string]string{"action": action, "otp": otp, "comment": "checked", "signature_base64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="})
		req := httptest.NewRequest("POST", "/api/v1/applications/"+strconv.Itoa(app.ID)+"/action", bytes.NewReader(payload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec.Code
	}
	readSignature := func(user models.User, role string) int {
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		req := httptest.NewRequest("GET", "/api/v1/applications/"+strconv.Itoa(app.ID)+"/"+role+"-signature", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec.Code
	}
	for _, action := range []string{"approve", "reject", "return_revision"} {
		if code := call(admin, action); code != 403 {
			t.Fatalf("admin %s at officer step: %d", action, code)
		}
	}
	if code := call(users[1], "approve"); code != 403 {
		t.Fatalf("director before officer: %d", code)
	}
	if code := call(users[0], "approve"); code != 200 {
		t.Fatalf("officer: %d", code)
	}
	directorToken, _ := helpers.GenerateToken(users[1].ID, users[1].Username, users[1].Role)
	detailReq := httptest.NewRequest("GET", "/api/v1/applications/"+strconv.Itoa(app.ID), nil)
	detailReq.Header.Set("Authorization", "Bearer "+directorToken)
	detailRec := httptest.NewRecorder()
	router.ServeHTTP(detailRec, detailReq)
	var detail map[string]json.RawMessage
	var data map[string]interface{}
	if detailRec.Code != 200 || json.Unmarshal(detailRec.Body.Bytes(), &detail) != nil || json.Unmarshal(detail["data"], &data) != nil || data["has_officer_signature"] != true {
		t.Fatalf("officer signature missing from application detail: %d", detailRec.Code)
	}
	signatureReq := httptest.NewRequest("GET", "/api/v1/applications/"+strconv.Itoa(app.ID)+"/officer-signature", nil)
	signatureReq.Header.Set("Authorization", "Bearer "+directorToken)
	signatureRec := httptest.NewRecorder()
	router.ServeHTTP(signatureRec, signatureReq)
	if signatureRec.Code != 200 || signatureRec.Header().Get("Content-Type") != "image/png" || !bytes.HasPrefix(signatureRec.Body.Bytes(), []byte{0x89, 0x50, 0x4e, 0x47}) {
		t.Fatalf("saved signature image is not returned: %d", signatureRec.Code)
	}
	if code := readSignature(users[1], "officer"); code != 200 {
		t.Fatalf("market director cannot read officer signature: %d", code)
	}
	if code := readSignature(users[3], "officer"); code != 403 {
		t.Fatalf("unassigned director read signature: %d", code)
	}
	if code := call(admin, "approve"); code != 403 {
		t.Fatalf("admin at market step: %d", code)
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
	if code := readSignature(users[2], "officer"); code != 200 {
		t.Fatalf("bureau director cannot read officer signature: %d", code)
	}
	if code := readSignature(users[2], "market"); code != 200 {
		t.Fatalf("bureau director cannot read market signature: %d", code)
	}
	if code := call(users[1], "approve"); code != 403 {
		t.Fatalf("market director at bureau step: %d", code)
	}
	if code := call(admin, "approve"); code != 403 {
		t.Fatalf("admin at bureau step: %d", code)
	}
	if code := call(users[2], "approve"); code != 422 {
		t.Fatalf("bureau director without OTP: %d", code)
	}
	otp = "123456"
	if code := call(users[2], "approve"); code != 200 {
		t.Fatalf("bureau director: %d", code)
	}
	if code := readSignature(users[0], "bureau"); code != 200 {
		t.Fatalf("officer cannot read bureau signature: %d", code)
	}
	var saved models.Application
	if err := db.First(&saved, app.ID).Error; err != nil {
		t.Fatal(err)
	}
	if saved.ApplicationState != "approved" || len(saved.OfficerSignature) == 0 || len(saved.MarketDirectorSignature) == 0 || len(saved.BureauDirectorSignature) == 0 {
		t.Fatal("approval or reviewer signatures were not saved")
	}

	adminToken, _ := helpers.GenerateToken(admin.ID, admin.Username, admin.Role)
	for _, path := range []string{"/api/v1/applications", "/api/v1/applications/" + strconv.Itoa(app.ID) + "/signed-document", "/api/v1/applications/clear"} {
		req := httptest.NewRequest("POST", path, bytes.NewBufferString("{}"))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != 403 {
			t.Fatalf("admin mutation %s: %d", path, rec.Code)
		}
	}
	req := httptest.NewRequest("GET", "/api/v1/applications/"+strconv.Itoa(app.ID), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("admin read application: %d", rec.Code)
	}
	var notifications int64
	if err := db.Model(&models.Notification{}).Where("type = ? AND target_url = ?", "application_signed", "/applications/"+strconv.Itoa(app.ID)).Count(&notifications).Error; err != nil || notifications != 1 {
		t.Fatalf("notification count = %d, error %v", notifications, err)
	}
	token, _ := helpers.GenerateToken(users[0].ID, users[0].Username, users[0].Role)
	upload := func(file string) int {
		payload, _ := json.Marshal(map[string]string{"file": file})
		req := httptest.NewRequest("POST", "/api/v1/applications/"+strconv.Itoa(app.ID)+"/signed-document", bytes.NewReader(payload))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec.Code
	}
	if code := upload("dGVzdA=="); code != 422 {
		t.Fatalf("invalid document: %d", code)
	}
	if code := upload("JVBERi0xLjQKZXhhbXBsZQ=="); code != 200 {
		t.Fatalf("valid document: %d", code)
	}
	getDocument := func(user models.User) *httptest.ResponseRecorder {
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		req := httptest.NewRequest("GET", "/api/v1/applications/"+strconv.Itoa(app.ID)+"/signed-document", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}
	document := getDocument(users[0])
	if document.Code != 200 || document.Header().Get("Content-Type") != "application/pdf" || !bytes.Equal(document.Body.Bytes(), []byte("%PDF-1.4\nexample")) {
		t.Fatalf("signed document did not match uploaded file: %d", document.Code)
	}
	if code := getDocument(users[3]).Code; code != 403 {
		t.Fatalf("unassigned user retrieved signed document: %d", code)
	}
	if code := upload("JVBERi0xLjQKZXhhbXBsZQ=="); code != 409 {
		t.Fatalf("repeat upload: %d", code)
	}
	if err := db.First(&saved, app.ID).Error; err != nil || saved.ApplicationState != "completed" || len(saved.SignedDocument) == 0 {
		t.Fatalf("document completion failed: %v", err)
	}
}
