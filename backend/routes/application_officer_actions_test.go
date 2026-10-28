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
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestOfficerReviewActionsAPI(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	databases.DB = db
	if err = databases.AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	market := models.Market{ID: 1, Slug: "officer-test", Code: "OT", Name: "Test", Status: "open"}
	if err = db.Create(&market).Error; err != nil {
		t.Fatal(err)
	}
	staff := models.User{Username: "reviewer", Password: "x", Role: "staff", Status: "A"}
	if err = db.Create(&staff).Error; err != nil {
		t.Fatal(err)
	}
	if err = db.Create(&models.MarketStaffAssignment{UserID: staff.ID, MarketID: market.ID}).Error; err != nil {
		t.Fatal(err)
	}
	token, err := helpers.GenerateToken(staff.ID, staff.Username, staff.Role)
	if err != nil {
		t.Fatal(err)
	}
	router := SetupRouter()
	cases := []struct {
		action, note, state string
		code                int
	}{{"approve", "checked", "pending_market_director", 200}, {"return_revision", "missing document", "revision_required", 200}, {"reject", "not eligible", "rejected", 200}, {"return_revision", "", "pending_officer", 422}, {"reject", "", "pending_officer", 422}}
	for i, tc := range cases {
		t.Run(tc.action+strconv.Itoa(i), func(t *testing.T) {
			app := models.Application{ApplicationNo: "TEST-OFFICER-" + strconv.Itoa(i), ApplicationType: "stall_lease", StallIdent: "TEST-" + strconv.Itoa(i), MarketID: &market.ID, ApplicantName: "Test", ApplicantPhone: "1", ShopName: "Test", ApplicationState: "pending_officer", Status: "A"}
			if err := db.Create(&app).Error; err != nil {
				t.Fatal(err)
			}
			payload, _ := json.Marshal(map[string]string{"action": tc.action, "comment": tc.note})
			req := httptest.NewRequest(http.MethodPost, "/api/v1/applications/"+strconv.Itoa(app.ID)+"/action", bytes.NewReader(payload))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			if res.Code != tc.code {
				t.Fatalf("status %d: %s", res.Code, res.Body.String())
			}
			if err := db.First(&app, app.ID).Error; err != nil {
				t.Fatal(err)
			}
			if app.ApplicationState != tc.state {
				t.Fatalf("state %q want %q", app.ApplicationState, tc.state)
			}
			if tc.note != "" && app.OfficerComment != tc.note {
				t.Fatalf("comment %q", app.OfficerComment)
			}
		})
	}
}
