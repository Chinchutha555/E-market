package handlers

import (
	"emarket-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func OfficerSignature(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		var app models.Application
		if err := db.Select("officer_signature", "officer_signature_type").Where("id = ? AND status = ?", id, "A").First(&app).Error; err != nil || len(app.OfficerSignature) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Header("Cache-Control", "private, no-store")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Data(http.StatusOK, app.OfficerSignatureType, app.OfficerSignature)
	}
}
