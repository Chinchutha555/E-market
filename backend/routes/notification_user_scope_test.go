package routes

import (
	"bytes"
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

func TestNotificationsAreScopedPerUser(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	databases.DB = db
	if err := databases.AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	market := models.Market{Slug: "notify-scope", Code: "N", Name: "Scoped market", Status: "open"}
	if err := db.Create(&market).Error; err != nil {
		t.Fatal(err)
	}
	users := []models.User{
		{Username: "notify-staff-one", Role: "staff", Status: "A"},
		{Username: "notify-staff-two", Role: "staff", Status: "A"},
		{Username: "notify-director", Role: "market_director", Status: "A"},
		{Username: "notify-tenant", Role: "tenant", Status: "A"},
	}
	if err := db.Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	for _, user := range users[:3] {
		if err := db.Create(&models.MarketStaffAssignment{UserID: user.ID, MarketID: market.ID}).Error; err != nil {
			t.Fatal(err)
		}
	}
	shared := models.Notification{RecipientRole: "officer", MarketName: market.Name, Type: "test", Title: "Staff notice"}
	personal := models.Notification{RecipientRole: "market_director", RecipientUserID: &users[2].ID, Type: "test", Title: "Director notice"}
	tenant := models.Notification{RecipientRole: "tenant", RecipientUserID: &users[3].ID, Type: "test", Title: "Tenant notice"}
	if err := db.Create(&[]models.Notification{shared, personal, tenant}).Error; err != nil {
		t.Fatal(err)
	}
	router := SetupRouter()
	call := func(user models.User, method, path string) *httptest.ResponseRecorder {
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		req := httptest.NewRequest(method, path, bytes.NewBufferString("{}"))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}
	list := func(user models.User) []models.Notification {
		rec := call(user, "GET", "/api/v1/notifications?role=officer")
		if rec.Code != 200 {
			t.Fatalf("list %s: %d %s", user.Username, rec.Code, rec.Body.String())
		}
		var result struct{ Data []models.Notification }
		if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		return result.Data
	}
	first, second, director, tenantList := list(users[0]), list(users[1]), list(users[2]), list(users[3])
	if len(first) != 1 || len(second) != 1 || len(director) != 1 || len(tenantList) != 1 {
		t.Fatalf("unexpected notification counts: %d %d %d %d", len(first), len(second), len(director), len(tenantList))
	}
	if first[0].ID == second[0].ID || first[0].RecipientUserID == nil || *first[0].RecipientUserID != users[0].ID || director[0].Title != "Director notice" {
		t.Fatal("notifications leaked across users")
	}
	if code := call(users[2], "POST", fmt.Sprintf("/api/v1/notifications/%d/read", first[0].ID)).Code; code != 404 {
		t.Fatalf("director read staff notice: %d", code)
	}
	if code := call(users[0], "POST", fmt.Sprintf("/api/v1/notifications/%d/read", first[0].ID)).Code; code != 200 {
		t.Fatalf("staff read own notice: %d", code)
	}
	if !list(users[0])[0].IsRead || list(users[1])[0].IsRead {
		t.Fatal("read state was shared across accounts")
	}
	if len(list(users[0])) != 1 {
		t.Fatal("duplicate notification copy")
	}
}
