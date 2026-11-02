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
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestBureauRejectRequiresOTPAndClosesApplication(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open("file:bureau-reject-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	databases.DB = db
	if err = databases.AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	tenant := models.User{Username: "bureau-test-tenant", Password: "x", FullName: "Tenant", Role: "tenant", Status: "A"}
	bureau := models.User{Username: "bureau-test-director", Password: "x", FullName: "Director", Role: "bureau_director", Status: "A"}
	if err = db.Create(&tenant).Error; err != nil {
		t.Fatal(err)
	}
	if err = db.Create(&bureau).Error; err != nil {
		t.Fatal(err)
	}
	app := models.Application{ApplicationNo: "BUREAU-REJECT", ApplicationType: "stall_lease", StallIdent: "S-BUREAU", ApplicantName: "Tenant", ApplicantPhone: "123", ShopName: "Shop", ApplicationState: "pending_bureau_director", Status: "A", CreatedBy: &tenant.ID}
	if err = db.Create(&app).Error; err != nil {
		t.Fatal(err)
	}
	router := SetupRouter()
	call := func(user models.User, path string, body interface{}) *httptest.ResponseRecorder {
		data, _ := json.Marshal(body)
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(data))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		return res
	}
	path := "/api/v1/applications/" + strconv.Itoa(app.ID)
	payload := map[string]interface{}{"action": "reject", "bureau_director_review": map[string]interface{}{"decision": "reject"}}
	if res := call(bureau, path+"/action", payload); res.Code != 422 {
		t.Fatalf("missing OTP: %d %s", res.Code, res.Body.String())
	}
	payload["otp"] = "000000"
	if res := call(bureau, path+"/action", payload); res.Code != 422 {
		t.Fatalf("wrong OTP: %d", res.Code)
	}
	payload["otp"] = "123456"
	if res := call(bureau, path+"/action", payload); res.Code != 200 {
		t.Fatalf("reject: %d %s", res.Code, res.Body.String())
	}
	var saved models.Application
	if err = db.First(&saved, app.ID).Error; err != nil {
		t.Fatal(err)
	}
	if saved.ApplicationState != "rejected" || len(saved.BureauDirectorSignature) != 0 || saved.BureauDirectorReviewJSON == "" {
		t.Fatalf("unexpected state: %#v", saved)
	}
	body := map[string]interface{}{"applicant_name": "Tenant", "applicant_phone": "123", "shop_name": "Shop"}
	if res := call(tenant, path+"/resubmit", body); res.Code != 409 {
		t.Fatalf("closed application resubmitted: %d %s", res.Code, res.Body.String())
	}
}
