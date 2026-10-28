package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"emarket-api/middlewares"
	"emarket-api/models"
	"emarket-api/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Resubmit updates only the tenant's returned application, without changing its booking.
func Resubmit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		var req struct {
			ApplicantName      string `json:"applicant_name"`
			ApplicantPhone     string `json:"applicant_phone"`
			ApplicantIDCard    string `json:"applicant_id_card"`
			ApplicantEmail     string `json:"applicant_email"`
			ApplicantAddress   string `json:"applicant_address"`
			ApplicantAge       *int   `json:"applicant_age"`
			ShopName           string `json:"shop_name"`
			ProductCategory    string `json:"product_category"`
			ProductType        string `json:"product_type"`
			Details            string `json:"details"`
			RequestedStartDate string `json:"requested_start_date"`
			LeaseTermMonths    int    `json:"lease_term_months"`
		}
		if err != nil || c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.ApplicantName) == "" || strings.TrimSpace(req.ApplicantPhone) == "" || strings.TrimSpace(req.ShopName) == "" {
			c.JSON(http.StatusUnprocessableEntity, responses.FailRespone(responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุชื่อ เบอร์โทรศัพท์ และชื่อร้าน")))
			return
		}

		var app models.Application
		if err := db.Where("id = ?", id).First(&app).Error; err != nil {
			c.JSON(http.StatusNotFound, responses.FailRespone(responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลคำร้อง")))
			return
		}

		userID := middlewares.GetUserID(c)
		if app.CreatedBy != nil && *app.CreatedBy != 0 && *app.CreatedBy != userID {
			c.JSON(http.StatusForbidden, responses.FailRespone(responses.NewAppErr(http.StatusForbidden, "สิทธิ์ของคุณไม่เพียงพอสำหรับการดำเนินการนี้")))
			return
		}

		if app.ApplicationState != "revision_required" && app.ApplicationState != "rejected" {
			c.JSON(http.StatusConflict, responses.FailRespone(responses.NewAppErr(http.StatusConflict, "คำร้องไม่อยู่ในสถานะที่สามารถส่งแก้ไขได้")))
			return
		}

		productCat := strings.TrimSpace(req.ProductCategory)
		if productCat == "" {
			productCat = strings.TrimSpace(req.ProductType)
		}
		if productCat == "" {
			productCat = "สินค้าทั่วไป"
		}

		updates := map[string]interface{}{
			"applicant_name":    strings.TrimSpace(req.ApplicantName),
			"applicant_phone":   strings.TrimSpace(req.ApplicantPhone),
			"applicant_id_card": strings.TrimSpace(req.ApplicantIDCard),
			"applicant_email":   strings.TrimSpace(req.ApplicantEmail),
			"applicant_address": strings.TrimSpace(req.ApplicantAddress),
			"shop_name":         strings.TrimSpace(req.ShopName),
			"product_category":  productCat,
			"details":           strings.TrimSpace(req.Details),
			"application_state": "pending_officer",
			"updated_by":        userID,
			"updated_at":        time.Now(),
		}
		if req.ApplicantAge != nil {
			updates["applicant_age"] = *req.ApplicantAge
		}
		if req.RequestedStartDate != "" {
			updates["requested_start_date"] = req.RequestedStartDate
		}
		if req.LeaseTermMonths > 0 {
			updates["lease_term_months"] = req.LeaseTermMonths
		}

		if err := db.Model(&app).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewInternalServerError()))
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"id": id, "status": "pending_officer"}, http.StatusOK))
	}
}
