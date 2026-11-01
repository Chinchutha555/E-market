package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"emarket-api/middlewares"
	"emarket-api/models"
	"emarket-api/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MaintenanceHandler struct{ db *gorm.DB }

func NewMaintenanceHandler(db *gorm.DB) *MaintenanceHandler { return &MaintenanceHandler{db: db} }

func maintenanceResponse(record models.MaintenanceRecord) map[string]any {
	payload := map[string]any{}
	_ = json.Unmarshal([]byte(record.PayloadJSON), &payload)
	payload["id"] = strconv.FormatUint(uint64(record.ID), 10)
	payload["market_id"] = record.MarketID
	payload["stallIdent"] = record.StallIdent
	payload["stall_ident"] = record.StallIdent
	payload["status"] = record.Status
	payload["createdAt"] = record.CreatedAt
	payload["updatedAt"] = record.UpdatedAt
	return payload
}

func maintenanceScope(c *gin.Context, db *gorm.DB, marketID int) bool {
	if marketID <= 0 || !middlewares.CanManageMarket(db, middlewares.GetUserID(c), middlewares.GetUserRole(c), marketID) {
		c.AbortWithStatus(http.StatusForbidden)
		return false
	}
	return true
}

func (h *MaintenanceHandler) List(c *gin.Context) {
	marketID, _ := strconv.Atoi(strings.TrimSpace(c.Query("market_id")))
	if !maintenanceScope(c, h.db, marketID) {
		return
	}
	query := h.db.Where("market_id = ?", marketID)
	if stallIdent := strings.ToUpper(strings.TrimSpace(c.Query("stall_ident"))); stallIdent != "" {
		query = query.Where("stall_ident = ?", stallIdent)
	}
	var records []models.MaintenanceRecord
	if err := query.Order("created_at DESC, id DESC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewAppErr(500, "ไม่สามารถโหลดประวัติงานซ่อม")))
		return
	}
	result := make([]map[string]any, 0, len(records))
	for _, record := range records {
		result = append(result, maintenanceResponse(record))
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(result, http.StatusOK))
}

func readMaintenancePayload(c *gin.Context) (map[string]any, int, string, string, bool) {
	payload := map[string]any{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responses.FailRespone(responses.NewAppErr(422, "ข้อมูลรายการซ่อมไม่ถูกต้อง")))
		return nil, 0, "", "", false
	}
	marketID, _ := strconv.Atoi(strings.TrimSpace(toString(payload["market_id"])))
	stallIdent := strings.ToUpper(strings.TrimSpace(firstString(payload, "stall_ident", "stallIdent")))
	status := strings.ToLower(strings.TrimSpace(toString(payload["status"])))
	if status == "" {
		status = "pending"
	}
	if marketID <= 0 || stallIdent == "" {
		c.JSON(http.StatusUnprocessableEntity, responses.FailRespone(responses.NewAppErr(422, "กรุณาระบุตลาดและแผงค้า")))
		return nil, 0, "", "", false
	}
	return payload, marketID, stallIdent, status, true
}

func toString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.FormatInt(int64(typed), 10)
	case json.Number:
		return typed.String()
	default:
		return ""
	}
}

func firstString(payload map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := toString(payload[key]); value != "" {
			return value
		}
	}
	return ""
}

func (h *MaintenanceHandler) Create(c *gin.Context) {
	payload, marketID, stallIdent, status, ok := readMaintenancePayload(c)
	if !ok || !maintenanceScope(c, h.db, marketID) {
		return
	}
	encoded, _ := json.Marshal(payload)
	record := models.MaintenanceRecord{MarketID: marketID, StallIdent: stallIdent, Status: status, PayloadJSON: string(encoded), CreatedBy: middlewares.GetUserID(c)}
	if err := h.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewAppErr(500, "ไม่สามารถบันทึกรายการซ่อม")))
		return
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse(maintenanceResponse(record), http.StatusCreated))
}

func (h *MaintenanceHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var record models.MaintenanceRecord
	if h.db.First(&record, id).Error != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !maintenanceScope(c, h.db, record.MarketID) {
		return
	}
	payload, marketID, stallIdent, status, ok := readMaintenancePayload(c)
	if !ok || marketID != record.MarketID || stallIdent != record.StallIdent {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	encoded, _ := json.Marshal(payload)
	record.Status, record.PayloadJSON = status, string(encoded)
	if err := h.db.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewAppErr(500, "ไม่สามารถอัปเดตรายการซ่อม")))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(maintenanceResponse(record), http.StatusOK))
}

func (h *MaintenanceHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var record models.MaintenanceRecord
	if h.db.First(&record, id).Error != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !maintenanceScope(c, h.db, record.MarketID) {
		return
	}
	if err := h.db.Delete(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewAppErr(500, "ไม่สามารถลบรายการซ่อม")))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"deleted": true}, http.StatusOK))
}
