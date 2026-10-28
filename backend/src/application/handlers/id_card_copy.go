package handlers

import (
	"emarket-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// IDCardCopy is reachable only after the application's resource-access check.
func IDCardCopy(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		var app models.Application
		if err := db.Select("id_card_copy", "id_card_copy_type").Where("id = ? AND status = ?", id, "A").First(&app).Error; err != nil || len(app.IDCardCopy) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Header("Cache-Control", "private, no-store")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Content-Disposition", `attachment; filename="id-card-copy"`)
		c.Data(http.StatusOK, app.IDCardCopyType, app.IDCardCopy)
	}
}
