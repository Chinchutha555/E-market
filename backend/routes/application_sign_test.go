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

func TestApplicationSignatureOwnerAndPersistence(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open("file:application-sign?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	databases.DB = db
	if err := databases.AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	users := []models.User{{Username: "sign-owner", Password: "x", Role: "tenant", Status: "A"}, {Username: "sign-other", Password: "x", Role: "tenant", Status: "A"}}
	if err := db.Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	owner := users[0].ID
	app := models.Application{ApplicationNo: "SIGN-TEST", ApplicationType: "stall_lease", StallIdent: "SIGN-01", ApplicantName: "Owner", ApplicantPhone: "1", ShopName: "Shop", ApplicationState: "pending_officer", Status: "A", CreatedBy: &owner}
	if err := db.Create(&app).Error; err != nil {
		t.Fatal(err)
	}
	router := SetupRouter()
	call := func(user models.User, method string) *httptest.ResponseRecorder {
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		req := httptest.NewRequest(method, "/api/v1/applications/"+strconv.Itoa(app.ID)+"/sign", bytes.NewBufferString(`{"name":"Owner"}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}
	if rec := call(users[1], "POST"); rec.Code != 403 {
		t.Fatalf("other tenant: %d", rec.Code)
	}
	if rec := call(users[0], "POST"); rec.Code != 200 {
		t.Fatalf("owner sign: %d %s", rec.Code, rec.Body.String())
	}
	var saved models.Application
	if err := db.First(&saved, app.ID).Error; err != nil {
		t.Fatal(err)
	}
	if saved.TenantSignature != "Owner" || saved.TenantSignedAt == nil {
		t.Fatalf("signature not persisted: %+v", saved)
	}
	if rec := call(users[0], "POST"); rec.Code != 409 {
		t.Fatalf("duplicate signature: %d", rec.Code)
	}
}
