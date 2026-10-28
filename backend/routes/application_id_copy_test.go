package routes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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

func TestApplicationIDCopyAccess(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	databases.DB = db
	if err := databases.AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	users := []models.User{{Username: "copy-owner", Password: "x", Role: "tenant", Status: "A"}, {Username: "copy-other", Password: "x", Role: "tenant", Status: "A"}}
	if err := db.Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	router := SetupRouter()
	call := func(user models.User, method, path string, body []byte) *httptest.ResponseRecorder {
		token, _ := helpers.GenerateToken(user.ID, user.Username, user.Role)
		req := httptest.NewRequest(method, path, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}
	file := []byte("%PDF-1.4\nexample")
	payload, _ := json.Marshal(map[string]interface{}{"application_type": "general", "applicant_name": "Test", "applicant_phone": "1", "shop_name": "Shop", "applicant_age": 30, "id_card_copy_base64": base64.StdEncoding.EncodeToString(file)})
	rec := call(users[0], "POST", "/api/v1/applications", payload)
	if rec.Code != 201 {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	var created struct {
		Data struct {
			ID           int  `json:"id"`
			ApplicantAge int  `json:"applicant_age"`
			HasCopy      bool `json:"has_id_card_copy"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.Data.ApplicantAge != 30 || !created.Data.HasCopy {
		t.Fatalf("missing age or copy: %+v", created.Data)
	}
	path := "/api/v1/applications/" + strconv.Itoa(created.Data.ID) + "/id-card-copy"
	if rec := call(users[1], "GET", path, nil); rec.Code != 403 {
		t.Fatalf("other tenant access: %d", rec.Code)
	}
	if rec := call(users[0], "GET", path, nil); rec.Code != 200 || !bytes.Equal(rec.Body.Bytes(), file) {
		t.Fatalf("owner access: %d", rec.Code)
	}
	bad, _ := json.Marshal(map[string]interface{}{"application_type": "general", "applicant_name": "Test", "applicant_phone": "1", "shop_name": "Shop", "id_card_copy_base64": base64.StdEncoding.EncodeToString([]byte("not a document"))})
	if rec := call(users[0], "POST", "/api/v1/applications", bad); rec.Code != 422 {
		t.Fatalf("invalid file: %d", rec.Code)
	}
}
