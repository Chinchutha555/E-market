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

func SignApplication(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		var body struct {
			Name string `json:"name"`
		}
		if err != nil || id <= 0 || c.ShouldBindJSON(&body) != nil || strings.TrimSpace(body.Name) == "" || len([]rune(strings.TrimSpace(body.Name))) > 255 {
			c.JSON(http.StatusUnprocessableEntity, responses.FailRespone(responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุชื่อผู้ลงนาม")))
			return
		}
		userID := middlewares.GetUserID(c)
		now := time.Now()
		result := db.Model(&models.Application{}).Where("id = ? AND created_by = ? AND status = ? AND application_type = ? AND (tenant_signature IS NULL OR tenant_signature = '') AND application_state NOT IN ?", id, userID, "A", "stall_lease", []string{"rejected", "completed"}).Updates(map[string]interface{}{"tenant_signature": strings.TrimSpace(body.Name), "tenant_signed_at": now, "updated_by": userID})
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewInternalServerError()))
			return
		}
		if result.RowsAffected != 1 {
			c.JSON(http.StatusConflict, responses.FailRespone(responses.NewAppErr(http.StatusConflict, "ไม่สามารถลงนามคำร้องนี้ได้ หรือคำร้องลงนามแล้ว")))
			return
		}
		c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"tenant_signature": strings.TrimSpace(body.Name), "tenant_signed_at": now}, http.StatusOK))
	}
}
