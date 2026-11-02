package handlers

import (
	"encoding/json"
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

		if app.ApplicationState == "rejected" && app.BureauDirectorReviewJSON != "" {
			var review models.BureauDirectorReview
			if json.Unmarshal([]byte(app.BureauDirectorReviewJSON), &review) != nil || review.Decision == "reject" {
				c.JSON(http.StatusConflict, responses.FailRespone(responses.NewAppErr(http.StatusConflict, "คำร้องนี้ปิดแล้ว")))
				return
			}
		}
		if app.ApplicationState != "revision_required" && app.ApplicationState != "rejected" && app.ApplicationState != "debt_pending" {
			c.JSON(http.StatusConflict, responses.FailRespone(responses.NewAppErr(http.StatusConflict, "คำร้องไม่อยู่ในสถานะที่สามารถส่งแก้ไขได้")))
			return
		}
		if app.ApplicationState == "debt_pending" && app.CreatedBy != nil && *app.CreatedBy > 0 {
			_ = db.Model(&models.Debt{}).Where("tenant_id = ? AND status = 'outstanding'", *app.CreatedBy).Update("status", "paid").Error
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
		if req.RequestedStartDate != "" && app.TradeOptionID == "" {
			updates["requested_start_date"] = req.RequestedStartDate
		}
		if req.LeaseTermMonths > 0 && app.TradeDuration == 0 {
			updates["lease_term_months"] = req.LeaseTermMonths
		}

		if err := db.Model(&app).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewInternalServerError()))
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"id": id, "status": "pending_officer"}, http.StatusOK))
	}
}

// PayDebtAndResubmit marks outstanding debt as paid and returns the application to pending_officer.
func PayDebtAndResubmit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.FailRespone(responses.NewAppErr(http.StatusBadRequest, "ID ไม่ถูกต้อง")))
			return
		}
		var app models.Application
		if err := db.Where("id = ?", id).First(&app).Error; err != nil {
			c.JSON(http.StatusNotFound, responses.FailRespone(responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลคำร้อง")))
			return
		}
		userID := middlewares.GetUserID(c)
		role := middlewares.GetUserRole(c)
		if role == "tenant" && app.CreatedBy != nil && *app.CreatedBy != 0 && *app.CreatedBy != userID {
			c.JSON(http.StatusForbidden, responses.FailRespone(responses.NewAppErr(http.StatusForbidden, "สิทธิ์ของคุณไม่เพียงพอสำหรับการดำเนินการนี้")))
			return
		}
		if app.ApplicationState != "debt_pending" {
			c.JSON(http.StatusConflict, responses.FailRespone(responses.NewAppErr(http.StatusConflict, "คำร้องไม่อยู่ในสถานะรอชำระหนี้ค้าง")))
			return
		}
		now := time.Now()
		if app.CreatedBy != nil && *app.CreatedBy > 0 {
			_ = db.Model(&models.Debt{}).Where("tenant_id = ? AND status = 'outstanding'", *app.CreatedBy).Updates(map[string]interface{}{
				"status":     "paid",
				"updated_at": now,
			}).Error
		}
		updates := map[string]interface{}{
			"application_state": "pending_officer",
			"review_note":       "ผู้เช่าชำระหนี้ค้างเรียบร้อยแล้ว ส่งให้เจ้าหน้าที่ตรวจสอบอีกครั้ง",
			"updated_by":        userID,
			"updated_at":        now,
		}
		if err := db.Model(&app).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewInternalServerError()))
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"id": id, "status": "pending_officer"}, http.StatusOK))
	}
}
