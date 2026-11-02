package handlers

import (
	"emarket-api/middlewares"
	"emarket-api/models"
	"emarket-api/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }
func (h *Handler) List(c *gin.Context) {
	role, userID := middlewares.GetUserRole(c), middlewares.GetUserID(c)
	q := h.db.Model(&models.Debt{})
	switch role {
	case "tenant":
		q = q.Where("tenant_id = ?", userID)
	case "staff", "market_director":
		ids, err := middlewares.StaffMarketIDs(h.db, userID)
		if err != nil {
			c.JSON(500, responses.FailRespone(responses.NewAppErr(500, "ไม่สามารถตรวจสอบสิทธิ์ตลาด")))
			return
		}
		if len(ids) == 0 {
			q = q.Where("1 = 0")
		} else {
			q = q.Where("market_id IN ?", ids)
		}
	case "admin", "bureau_director", "finance":
	default:
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if rawMarketID := strings.TrimSpace(c.Query("market_id")); rawMarketID != "" {
		marketID, err := strconv.Atoi(rawMarketID)
		if err != nil || marketID <= 0 {
			c.JSON(http.StatusBadRequest, responses.FailRespone(responses.NewAppErr(http.StatusBadRequest, "รหัสตลาดไม่ถูกต้อง")))
			return
		}
		q = q.Where("market_id = ?", marketID)
	}
	var debts []models.Debt
	if err := q.Order("id DESC").Find(&debts).Error; err != nil {
		c.JSON(500, responses.FailRespone(responses.NewAppErr(500, "ไม่สามารถโหลดรายการหนี้")))
		return
	}
	c.JSON(200, responses.SuccessResponse(debts, 200))
}
func (h *Handler) Create(c *gin.Context) {
	var req struct {
		TenantID     int    `json:"tenant_id"`
		MarketID     int    `json:"market_id"`
		Description  string `json:"description"`
		AmountSatang int64  `json:"amount_satang"`
	}
	if c.ShouldBindJSON(&req) != nil || req.TenantID <= 0 || req.MarketID <= 0 || req.AmountSatang <= 0 || len(strings.TrimSpace(req.Description)) == 0 || len([]rune(req.Description)) > 500 {
		c.JSON(422, responses.FailRespone(responses.NewAppErr(422, "กรุณาระบุผู้เช่า ตลาด รายละเอียด และยอดหนี้ให้ถูกต้อง")))
		return
	}
	if !middlewares.CanManageMarket(h.db, middlewares.GetUserID(c), middlewares.GetUserRole(c), req.MarketID) {
		c.AbortWithStatus(403)
		return
	}
	var tenant models.User
	if err := h.db.Where("id = ? AND role = ? AND status = ?", req.TenantID, "tenant", "A").First(&tenant).Error; err != nil {
		c.JSON(422, responses.FailRespone(responses.NewAppErr(422, "ไม่พบผู้เช่าที่ระบุ")))
		return
	}
	var market models.Market
	if err := h.db.Where("id = ? AND status != ?", req.MarketID, "D").First(&market).Error; err != nil {
		c.JSON(422, responses.FailRespone(responses.NewAppErr(422, "ไม่พบตลาดที่ระบุ")))
		return
	}
	debt := models.Debt{TenantID: req.TenantID, MarketID: req.MarketID, Description: strings.TrimSpace(req.Description), AmountSatang: req.AmountSatang, Status: "outstanding", CreatedBy: middlewares.GetUserID(c)}
	if err := h.db.Create(&debt).Error; err != nil {
		c.JSON(500, responses.FailRespone(responses.NewAppErr(500, "ไม่สามารถบันทึกหนี้")))
		return
	}
	c.JSON(201, responses.SuccessResponse(debt, 201))
}

func (h *Handler) Pay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, responses.FailRespone(responses.NewAppErr(http.StatusBadRequest, "รหัสรายการหนี้ไม่ถูกต้อง")))
		return
	}
	var debt models.Debt
	if err := h.db.Where("id = ?", id).First(&debt).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.FailRespone(responses.NewAppErr(http.StatusNotFound, "ไม่พบรายการหนี้")))
		return
	}
	role, userID := middlewares.GetUserRole(c), middlewares.GetUserID(c)
	if role == "tenant" && debt.TenantID != userID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if err := h.db.Model(&debt).Update("status", "paid").Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewAppErr(http.StatusInternalServerError, "ไม่สามารถบันทึกการชำระหนี้")))
		return
	}
	var remaining int64
	h.db.Model(&models.Debt{}).Where("tenant_id = ? AND status = 'outstanding'", debt.TenantID).Count(&remaining)
	if remaining == 0 {
		_ = h.db.Model(&models.Application{}).
			Where("created_by = ? AND application_state = 'debt_pending'", debt.TenantID).
			Updates(map[string]interface{}{
				"application_state": "pending_officer",
				"review_note":       "ผู้เช่าชำระหนี้ค้างเรียบร้อยแล้ว ส่งให้เจ้าหน้าที่ตรวจสอบอีกครั้ง",
			}).Error
	}
	debt.Status = "paid"
	c.JSON(http.StatusOK, responses.SuccessResponse(debt, http.StatusOK))
}
