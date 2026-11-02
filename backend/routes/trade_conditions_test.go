package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"emarket-api/config"
	"emarket-api/databases"
	"emarket-api/helpers"
	"emarket-api/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestTradeConditionsCustomTypes(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open("file:trade-conditions-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	databases.DB = db
	if err = databases.AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	m1 := models.Market{ID: 1, Slug: "t1", Code: "T1", Name: "One", Status: "open"}
	m2 := models.Market{ID: 2, Slug: "t2", Code: "T2", Name: "Two", Status: "open"}
	staff := models.User{ID: 31, Username: "trade-staff", Password: "x", FullName: "Staff", Role: "staff", Status: "A"}
	for _, v := range []interface{}{&m1, &m2, &staff, &models.MarketStaffAssignment{UserID: 31, MarketID: 1}} {
		if err := db.Create(v).Error; err != nil {
			t.Fatal(err)
		}
	}
	token, _ := helpers.GenerateToken(staff.ID, staff.Username, staff.Role)
	router := SetupRouter()
	request := func(method, path string, body interface{}) *httptest.ResponseRecorder {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(method, path, bytes.NewReader(data))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		return res
	}
	one := map[string]interface{}{"TradeType": "custom1", "Title": "Custom", "RequiresContract": false, "MinDuration": 1, "MaxDuration": 3, "Unit": "days"}
	two := map[string]interface{}{"TradeType": "custom2", "Title": "Another", "MinDuration": 2, "MaxDuration": 4, "Unit": "years", "StartDate": "2024-02-29", "EndDate": "2026-02-28"}
	one["Options"] = []interface{}{
		map[string]interface{}{"ID": "one", "Duration": 1, "Unit": "months", "StartDate": "2026-11-02", "EndDate": "2026-12-01"},
		map[string]interface{}{"ID": "three", "Duration": 3, "Unit": "months", "StartDate": "2026-10-01", "EndDate": "2026-12-31"},
	}
	payload := map[string]interface{}{"Conditions": []interface{}{one, two}}
	if res := request(http.MethodPut, "/api/v1/markets/1/trade-conditions", payload); res.Code != 200 {
		t.Fatalf("save: %d %s", res.Code, res.Body.String())
	}
	if res := request(http.MethodGet, "/api/v1/markets/1/trade-conditions", nil); res.Code != 200 || !bytes.Contains(res.Body.Bytes(), []byte("Custom")) {
		t.Fatalf("load: %d %s", res.Code, res.Body.String())
	}
	var saved models.MarketTradeCondition
	if err := db.Where("market_id = ? AND trade_type = ?", 1, "custom1").First(&saved).Error; err != nil {
		t.Fatal(err)
	}
	if len(saved.DurationOptions()) != 2 || saved.DurationOptions()[1].ID != "three" {
		t.Fatal("multiple options were not saved")
	}
	if saved.RequiresContract == nil || *saved.RequiresContract {
		t.Fatal("no-contract setting was not saved")
	}
	if res := request(http.MethodPut, "/api/v1/markets/2/trade-conditions", payload); res.Code != 403 {
		t.Fatalf("market scope: %d", res.Code)
	}
	if res := request(http.MethodPut, "/api/v1/markets/1/trade-conditions", map[string]interface{}{"Conditions": []interface{}{two}}); res.Code != 200 {
		t.Fatalf("remove: %d %s", res.Code, res.Body.String())
	}
	var count int64
	db.Model(&models.MarketTradeCondition{}).Where("market_id = ?", 1).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 remaining row, got %d", count)
	}
	singleDay := map[string]interface{}{"TradeType": "single-day", "Title": "Single day", "MinDuration": 1, "MaxDuration": 1, "Unit": "days", "StartDate": "2026-10-01", "EndDate": "2026-10-01"}
	if res := request(http.MethodPut, "/api/v1/markets/1/trade-conditions", map[string]interface{}{"Conditions": []interface{}{singleDay}}); res.Code != 200 {
		t.Fatalf("single-day dates rejected: %d %s", res.Code, res.Body.String())
	}
	invalidDates := map[string]interface{}{"TradeType": "bad-date", "Title": "Bad date", "MinDuration": 1, "MaxDuration": 2, "Unit": "days", "StartDate": "2026-02-28", "EndDate": "2026-02-01"}
	if res := request(http.MethodPut, "/api/v1/markets/1/trade-conditions", map[string]interface{}{"Conditions": []interface{}{invalidDates}}); res.Code != 400 {
		t.Fatalf("invalid dates accepted: %d", res.Code)
	}
	bad := map[string]interface{}{"TradeType": "other", "Title": "Another", "MinDuration": 1, "MaxDuration": 3, "Unit": "days"}
	if res := request(http.MethodPut, "/api/v1/markets/1/trade-conditions", map[string]interface{}{"Conditions": []interface{}{two, bad}}); res.Code != 400 {
		t.Fatalf("duplicate names: %d", res.Code)
	}
	db.Model(&models.MarketTradeCondition{}).Where("market_id = ?", 1).Count(&count)
	if count != 1 {
		t.Fatalf("invalid save changed rows: %d", count)
	}
	if res := request(http.MethodPut, "/api/v1/markets/1/trade-conditions", map[string]interface{}{"Conditions": []interface{}{}}); res.Code != 200 {
		t.Fatalf("empty list: %d %s", res.Code, res.Body.String())
	}
}
