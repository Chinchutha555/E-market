package middlewares

import (
	"emarket-api/models"
	"emarket-api/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// RequireApplicationReviewer enforces the role and market for the current step.
func RequireApplicationReviewer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		var app models.Application
		if err := db.Where("id = ? AND status = ?", id, "A").First(&app).Error; err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		role := GetUserRole(c)
		allowed := false
		switch app.ApplicationState {
		case "pending_officer":
			allowed = role == "admin" || (role == "staff" && (app.MarketID == nil || CanManageMarket(db, GetUserID(c), role, *app.MarketID)))
		case "pending_market_director":
			allowed = role == "admin" || (role == "market_director" && (app.MarketID == nil || CanManageMarket(db, GetUserID(c), role, *app.MarketID)))
		case "pending_bureau_director":
			allowed = role == "admin" || role == "bureau_director"
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, responses.FailRespone(responses.NewAppErr(http.StatusForbidden, "สิทธิ์ของคุณไม่เพียงพอสำหรับการดำเนินการนี้")))
			return
		}
		c.Next()
	}
}
