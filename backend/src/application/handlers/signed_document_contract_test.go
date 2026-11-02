package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"emarket-api/models"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestSignedDocumentWithoutContractNotifiesCompletion(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:signed-no-contract?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&models.Booking{}, &models.Application{}, &models.Notification{}); err != nil {
		t.Fatal(err)
	}
	noContract := false
	userID := 10
	booking := models.Booking{BookingNo: "NO-CONTRACT-BK", StallIdent: "A-1", ApplicantName: "Tenant", ApplicantPhone: "123", VisitDate: "2026-10-01", VisitTime: "10:00", RequiresContract: &noContract}
	if err = db.Create(&booking).Error; err != nil {
		t.Fatal(err)
	}
	app := models.Application{ApplicationNo: "NO-CONTRACT-APP", ApplicationType: "stall_lease", BookingID: &booking.ID, StallIdent: "A-1", ApplicantName: "Tenant", ApplicantPhone: "123", ShopName: "Shop", ApplicationState: "approved", Status: "A", CreatedBy: &userID}
	if err = db.Create(&app).Error; err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.POST("/:id", SignedDocumentUpload(db))
	body, _ := json.Marshal(map[string]string{"File": base64.StdEncoding.EncodeToString([]byte("%PDF-1.4\n%%EOF"))})
	req := httptest.NewRequest(http.MethodPost, "/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("upload: %d %s", res.Code, res.Body.String())
	}
	var notification models.Notification
	if err = db.First(&notification).Error; err != nil {
		t.Fatal(err)
	}
	if notification.Type != "application_completed" || notification.TargetURL != "/applications" {
		t.Fatalf("wrong notification: %#v", notification)
	}
}
