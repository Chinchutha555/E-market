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
	"testing"
)

func TestDebtAccessAndCreation(t *testing.T) {
	config.Init()
	db, e := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if e != nil {
		t.Fatal(e)
	}
	databases.DB = db
	if e = databases.AutoMigrate(db); e != nil {
		t.Fatal(e)
	}
	users := []models.User{{Username: "officer-debt", Password: "x", Role: "staff", Status: "A"}, {Username: "tenant-debt", Password: "x", Role: "tenant", Status: "A"}, {Username: "other-debt", Password: "x", Role: "tenant", Status: "A"}}
	if e = db.Create(&users).Error; e != nil {
		t.Fatal(e)
	}
	markets := []models.Market{{ID: 1, Slug: "debt-one", Code: "D1", Name: "One", Status: "open"}, {ID: 2, Slug: "debt-two", Code: "D2", Name: "Two", Status: "open"}}
	if e = db.Create(&markets).Error; e != nil {
		t.Fatal(e)
	}
	if e = db.Create(&models.MarketStaffAssignment{UserID: users[0].ID, MarketID: 1}).Error; e != nil {
		t.Fatal(e)
	}
	router := SetupRouter()
	call := func(user models.User, method, path string, body any) *httptest.ResponseRecorder {
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}
	valid := map[string]any{"tenant_id": users[1].ID, "market_id": 1, "description": "ค่าเช่า", "amount_satang": 12345}
	if rec := call(users[0], "POST", "/api/v1/debts", valid); rec.Code != 201 {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	invalid := map[string]any{"tenant_id": users[1].ID, "market_id": 2, "description": "ค่าเช่า", "amount_satang": 12345}
	if rec := call(users[0], "POST", "/api/v1/debts", invalid); rec.Code != 403 {
		t.Fatalf("foreign market: %d", rec.Code)
	}
	if rec := call(users[1], "POST", "/api/v1/debts", valid); rec.Code != 403 {
		t.Fatalf("tenant create: %d", rec.Code)
	}
	for _, tc := range []struct {
		user  models.User
		count int
	}{{users[0], 1}, {users[1], 1}, {users[2], 0}} {
		rec := call(tc.user, "GET", "/api/v1/debts", nil)
		var result struct {
			Data []models.Debt `json:"data"`
		}
		if rec.Code != 200 {
			t.Fatalf("list %d", rec.Code)
		}
		if e = json.Unmarshal(rec.Body.Bytes(), &result); e != nil {
			t.Fatal(e)
		}
		if len(result.Data) != tc.count {
			t.Fatalf("%s saw %d debts", tc.user.Username, len(result.Data))
		}
	}
}
