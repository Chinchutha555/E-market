package routes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"emarket-api/config"
	"emarket-api/databases"
	"emarket-api/helpers"
	"emarket-api/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestContractReviewFlow(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	databases.DB = db
	if err := databases.AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	market := models.Market{Slug: "contract-flow", Code: "CF", Name: "Contract flow", Status: "open"}
	if err := db.Create(&market).Error; err != nil {
		t.Fatal(err)
	}
	users := []models.User{
		{Username: "contract-tenant", Role: "tenant", Status: "A"},
		{Username: "contract-staff", Role: "staff", Status: "A"},
		{Username: "contract-market", Role: "market_director", Status: "A"},
		{Username: "contract-bureau", Role: "bureau_director", Status: "A"},
		{Username: "contract-other", Role: "tenant", Status: "A"},
		{Username: "contract-admin", Role: "admin", Status: "A"},
	}
	if err := db.Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	for _, user := range users[1:3] {
		if err := db.Create(&models.MarketStaffAssignment{UserID: user.ID, MarketID: market.ID}).Error; err != nil {
			t.Fatal(err)
		}
	}
	stall := models.Stall{Ident: "FLOW-1", Name: "Flow", MarketID: market.ID, StallStatus: "reserved", Status: "A"}
	if err := db.Create(&stall).Error; err != nil {
		t.Fatal(err)
	}
	app := models.Application{ApplicationNo: "APP-CONTRACT-FLOW", ApplicationType: "stall_lease", StallIdent: stall.Ident, StallID: &stall.ID, MarketID: &market.ID, CreatedBy: &users[0].ID, ApplicantName: "Tenant", ApplicantPhone: "1", ShopName: "Shop", ApplicationState: "completed", Status: "A"}
	if err := db.Create(&app).Error; err != nil {
		t.Fatal(err)
	}
	router := SetupRouter()
	call := func(user models.User, method, path string, payload any) *httptest.ResponseRecorder {
		var body []byte
		if payload != nil {
			body, _ = json.Marshal(payload)
		}
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		req := httptest.NewRequest(method, path, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}
	signature := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="
	document := base64.StdEncoding.EncodeToString([]byte("%PDF-1.4\ncontract"))
	create := map[string]interface{}{"application_id": app.ID, "stall_id": stall.ID, "stall_ident": stall.Ident, "tenant_name": "Tenant", "tenant_phone": "0812345678", "tenant_id_card": "1234567890123", "shop_name": "Shop", "start_date": "2026-11-02", "end_date": "2027-11-02", "monthly_rental": 2000, "security_deposit": 2000, "signature_base64": signature, "document_base64": document}
	if rec := call(users[5], "POST", "/api/v1/contracts", create); rec.Code != 403 {
		t.Fatalf("admin create: %d", rec.Code)
	}
	missing := map[string]interface{}{}
	for k, v := range create {
		missing[k] = v
	}
	delete(missing, "document_base64")
	if rec := call(users[0], "POST", "/api/v1/contracts", missing); rec.Code != 422 {
		t.Fatalf("missing document: %d %s", rec.Code, rec.Body.String())
	}
	rec := call(users[0], "POST", "/api/v1/contracts", create)
	if rec.Code != 201 {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	var result struct {
		Data struct {
			ID     int
			Status string
		}
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Data.Status != "pending" {
		t.Fatalf("initial status: %s", result.Data.Status)
	}
	id := result.Data.ID
	path := fmt.Sprintf("/api/v1/contracts/%d", id)
	approve := map[string]string{"action": "approve", "signature": signature}
	if rec := call(users[2], "POST", path+"/action", approve); rec.Code != 403 {
		t.Fatalf("out of order: %d", rec.Code)
	}
	if rec := call(users[1], "POST", path+"/action", map[string]string{"action": "approve"}); rec.Code != 422 {
		t.Fatalf("missing signature: %d", rec.Code)
	}
	if rec := call(users[1], "POST", path+"/action", map[string]string{"action": "return_revision", "reason": "กรุณาแก้ไข"}); rec.Code != 200 {
		t.Fatalf("return: %d %s", rec.Code, rec.Body.String())
	}
	resubmit := map[string]interface{}{"TenantName": "Tenant", "TenantPhone": "0812345678", "TenantIDCard": "1234567890123", "ShopName": "Shop revised", "StartDate": "2026-11-02", "EndDate": "2027-11-02", "MonthlyRental": 2000, "SecurityDeposit": 2000, "Signature": signature, "Document": document}
	if rec := call(users[4], "POST", path+"/resubmit", resubmit); rec.Code != 403 {
		t.Fatalf("other tenant resubmit: %d", rec.Code)
	}
	if rec := call(users[0], "POST", path+"/resubmit", resubmit); rec.Code != 200 {
		t.Fatalf("resubmit: %d %s", rec.Code, rec.Body.String())
	}
	if rec := call(users[1], "POST", path+"/action", approve); rec.Code != 200 {
		t.Fatalf("officer: %d %s", rec.Code, rec.Body.String())
	}
	if rec := call(users[2], "GET", path+"/assets/officer", nil); rec.Code != 200 || rec.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("officer signature: %d", rec.Code)
	}
	if rec := call(users[2], "POST", path+"/action", map[string]string{"action": "return_revision", "reason": "แก้ไขเงื่อนไข"}); rec.Code != 200 {
		t.Fatalf("market return: %d", rec.Code)
	}
	if rec := call(users[0], "POST", path+"/resubmit", resubmit); rec.Code != 200 {
		t.Fatalf("market revision resubmit: %d %s", rec.Code, rec.Body.String())
	}
	if rec := call(users[1], "POST", path+"/action", approve); rec.Code != 200 {
		t.Fatalf("officer after revision: %d", rec.Code)
	}
	if rec := call(users[2], "POST", path+"/action", approve); rec.Code != 200 {
		t.Fatalf("market: %d %s", rec.Code, rec.Body.String())
	}
	if rec := call(users[3], "POST", path+"/action", approve); rec.Code != 422 {
		t.Fatalf("missing otp: %d", rec.Code)
	}
	if rec := call(users[3], "POST", path+"/action", map[string]string{"action": "approve", "signature": signature, "otp": "123456"}); rec.Code != 200 {
		t.Fatalf("bureau: %d %s", rec.Code, rec.Body.String())
	}
	var saved models.Contract
	if err := db.First(&saved, id).Error; err != nil {
		t.Fatal(err)
	}
	if saved.ContractState != "active" || len(saved.TenantSignature) == 0 || len(saved.OfficerSignature) == 0 || len(saved.MarketDirectorSignature) == 0 || len(saved.BureauDirectorSignature) == 0 {
		t.Fatal("contract signatures or final status missing")
	}
	if rec := call(users[0], "GET", path+"/assets/document", nil); rec.Code != 200 || !bytes.Equal(rec.Body.Bytes(), []byte("%PDF-1.4\ncontract")) {
		t.Fatalf("document: %d", rec.Code)
	}
	var count int64
	if err := db.Model(&models.Notification{}).Where("recipient_user_id = ? AND type = ? AND status_code = ?", users[0].ID, "contract_workflow", "active").Count(&count).Error; err != nil || count != 1 {
		t.Fatalf("tenant notification: %d %v", count, err)
	}
}
