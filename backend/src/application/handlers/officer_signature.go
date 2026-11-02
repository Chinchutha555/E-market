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

func DirectorSignature(db *gorm.DB, director string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		var app models.Application
		if err := db.First(&app, id).Error; err != nil {
			c.AbortWithStatus(404)
			return
		}
		data, kind := app.MarketDirectorSignature, app.MarketDirectorSignatureType
		if director == "bureau" {
			data, kind = app.BureauDirectorSignature, app.BureauDirectorSignatureType
		}
		if len(data) == 0 {
			c.AbortWithStatus(404)
			return
		}
		c.Header("Cache-Control", "private, no-store")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Data(200, kind, data)
	}
}
