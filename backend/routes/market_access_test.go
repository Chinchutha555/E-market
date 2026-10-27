package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"emarket-api/config"
	"emarket-api/databases"
	"emarket-api/helpers"
	"emarket-api/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func accessTestRequest(t *testing.T, router http.Handler, method, path, token string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("X-Active-Role", "admin") // must have no effect
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func TestStaffMarketScope(t *testing.T) {
	config.Init()
	db, err := gorm.Open(sqlite.Open("file:market-access-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	databases.DB = db
	if err := databases.AutoMigrate(db); err != nil {
		t.Fatal(err)
	}

	market1 := models.Market{ID: 1, Slug: "one", Code: "MKT-01", Name: "ตลาดหนึ่ง", Status: "open"}
	market2 := models.Market{ID: 2, Slug: "two", Code: "MKT-02", Name: "ตลาดสอง", Status: "open"}
	staff := models.User{ID: 21, Username: "staff-test", Password: "x", FullName: "เจ้าหน้าที่ทดสอบ", Role: "staff", Status: "A"}
	tenant := models.User{ID: 22, Username: "tenant-test", Password: "x", FullName: "ผู้เช่าทดสอบ", Role: "tenant", Status: "A"}
	for _, value := range []interface{}{&market1, &market2, &staff, &tenant, &models.MarketStaffAssignment{UserID: staff.ID, MarketID: market1.ID}} {
		if err := db.Create(value).Error; err != nil {
			t.Fatal(err)
		}
	}
	b1, b2 := models.Booking{BookingNo: "BK-MARKET-ONE", StallIdent: "ONE-01", MarketID: &market1.ID, ApplicantName: "A", ApplicantPhone: "1", VisitDate: "2026-01-01", VisitTime: "10:00", Status: "A"}, models.Booking{BookingNo: "BK-MARKET-TWO", StallIdent: "TWO-01", MarketID: &market2.ID, ApplicantName: "B", ApplicantPhone: "2", VisitDate: "2026-01-01", VisitTime: "10:00", Status: "A"}
	if err := db.Create(&b1).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&b2).Error; err != nil {
		t.Fatal(err)
	}

	staffToken, _ := helpers.GenerateToken(staff.ID, staff.Username, staff.Role)
	tenantToken, _ := helpers.GenerateToken(tenant.ID, tenant.Username, tenant.Role)
	router := SetupRouter()

	if got := accessTestRequest(t, router, http.MethodGet, "/api/v1/markets/1/gis-config/draft", staffToken); got.Code != http.StatusOK {
		t.Fatalf("assigned market should be accessible, got %d: %s", got.Code, got.Body.String())
	}
	if got := accessTestRequest(t, router, http.MethodGet, "/api/v1/markets/2/gis-config/draft", staffToken); got.Code != http.StatusForbidden {
		t.Fatalf("other market must be forbidden, got %d: %s", got.Code, got.Body.String())
	}
	if got := accessTestRequest(t, router, http.MethodGet, "/api/v1/markets/1/gis-config/draft", tenantToken); got.Code != http.StatusForbidden {
		t.Fatalf("X-Active-Role must not elevate a tenant, got %d", got.Code)
	}
	bookings := accessTestRequest(t, router, http.MethodGet, "/api/v1/bookings", staffToken)
	if bookings.Code != http.StatusOK || !strings.Contains(bookings.Body.String(), "BK-MARKET-ONE") || strings.Contains(bookings.Body.String(), "BK-MARKET-TWO") {
		t.Fatalf("booking list was not scoped to assigned market: %s", bookings.Body.String())
	}
}
