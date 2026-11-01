package usecases

import (
	"encoding/json"
	"strconv"
	"testing"

	"emarket-api/models"
	"emarket-api/requests"
	marketRepo "emarket-api/src/market/repositories"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newUtilityTestUseCase(t *testing.T) (*models.Market, func(string) (json.RawMessage, error), func(string, string) error) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Market{}, &models.MarketUtilitySetting{}); err != nil {
		t.Fatal(err)
	}
	market := &models.Market{Name: "ตลาดทดสอบ", Slug: "test-market"}
	if err := db.Create(market).Error; err != nil {
		t.Fatal(err)
	}
	uc := NewUseCase(marketRepo.NewRepository(db))
	get := func(id string) (json.RawMessage, error) {
		res, err := uc.GetUtilitySettings(id)
		if err != nil {
			return nil, err
		}
		return res.Settings, nil
	}
	save := func(id, body string) error {
		_, err := uc.SaveUtilitySettings(id, 1, requests.UpdateMarketUtilitySettingsRequest{Settings: json.RawMessage(body)})
		return err
	}
	return market, get, save
}

func TestUtilitySettingsRoundTrip(t *testing.T) {
	market, get, save := newUtilityTestUseCase(t)
	id := strconv.Itoa(market.ID)

	got, err := get(id)
	if err != nil || string(got) != "null" {
		t.Fatalf("ตลาดที่ยังไม่เคยตั้งค่า ต้องได้ null: %s %v", got, err)
	}

	if err := save(id, `{"scope":"project","defaults":{"electric":{"mode":"meter","rate":8}}}`); err != nil {
		t.Fatal(err)
	}
	// บันทึกซ้ำต้องแทนที่ ไม่สร้างแถวใหม่
	if err := save(id, `{"scope":"stall","defaults":{"electric":{"mode":"flat","flat":150}}}`); err != nil {
		t.Fatal(err)
	}
	got, err = get("test-market")
	if err != nil {
		t.Fatal(err)
	}
	var parsed struct {
		Scope string `json:"scope"`
	}
	if err := json.Unmarshal(got, &parsed); err != nil || parsed.Scope != "stall" {
		t.Fatalf("ต้องได้ค่าที่บันทึกล่าสุด: %s %v", got, err)
	}
}

func TestUtilitySettingsRejectsInvalidJSON(t *testing.T) {
	market, _, save := newUtilityTestUseCase(t)
	if err := save(strconv.Itoa(market.ID), `{not json`); err == nil {
		t.Fatal("JSON ที่ไม่ถูกต้องต้องถูกปฏิเสธ")
	}
	if err := save("999999", `{}`); err == nil {
		t.Fatal("ตลาดที่ไม่มีอยู่ต้องถูกปฏิเสธ")
	}
}
