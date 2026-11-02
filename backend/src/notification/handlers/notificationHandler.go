package handlers

import (
	"net/http"
	"strconv"

	"emarket-api/middlewares"
	"emarket-api/models"
	"emarket-api/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) scopedQuery(c *gin.Context) *gorm.DB {
	role := middlewares.GetUserRole(c)
	query := h.db.Model(&models.Notification{}).Where("recipient_user_id = ?", middlewares.GetUserID(c))
	switch role {
	case "staff":
		return query.Where("recipient_role IN ?", []string{"officer", "staff"})
	case "tenant", "market_director", "bureau_director", "admin":
		return query.Where("recipient_role = ?", role)
	default:
		return query.Where("1 = 0")
	}
}

// Copy market notices for each assigned staff member to keep read state private.
func (h *Handler) materializeStaffNotifications(c *gin.Context) error {
	if middlewares.GetUserRole(c) != "staff" {
		return nil
	}
	userID := middlewares.GetUserID(c)
	marketIDs, err := middlewares.StaffMarketIDs(h.db, userID)
	if err != nil || len(marketIDs) == 0 {
		return err
	}
	var sources []models.Notification
	err = h.db.Where("recipient_user_id IS NULL AND source_notification_id IS NULL AND recipient_role IN ?", []string{"officer", "staff"}).
		Where("booking_id IN (SELECT id FROM booking WHERE market_id IN ?) OR market_name IN (SELECT name FROM market WHERE id IN ? AND name IN (SELECT name FROM market GROUP BY name HAVING COUNT(*) = 1))", marketIDs, marketIDs).
		Where("NOT EXISTS (SELECT 1 FROM notification AS recipient WHERE recipient.source_notification_id = notification.id AND recipient.recipient_user_id = ?)", userID).Find(&sources).Error
	if err != nil {
		return err
	}
	for _, source := range sources {
		copy := source
		copy.ID = 0
		copy.SourceNotificationID = &source.ID
		copy.RecipientUserID = &userID
		copy.IsRead = false
		if err := h.db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "source_notification_id"}, {Name: "recipient_user_id"}}, DoNothing: true}).Create(&copy).Error; err != nil {
			return err
		}
	}
	return nil
}
func (h *Handler) GetNotifications(c *gin.Context) {
	if err := h.materializeStaffNotifications(c); err != nil {
		c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewInternalServerError()))
		return
	}
	var items []models.Notification
	if err := h.scopedQuery(c).Order("id DESC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewInternalServerError()))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(items, http.StatusOK))
}

func (h *Handler) MarkRead(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.FailRespone(responses.NewAppErr(http.StatusBadRequest, "ID การแจ้งเตือนไม่ถูกต้อง")))
		return
	}
	result := h.scopedQuery(c).Where("id = ?", id).Update("is_read", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewInternalServerError()))
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, responses.FailRespone(responses.NewAppErr(http.StatusNotFound, "ไม่พบการแจ้งเตือน")))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"id": id, "is_read": true}, http.StatusOK))
}

func (h *Handler) MarkAllRead(c *gin.Context) {
	if err := h.scopedQuery(c).Where("is_read = ?", false).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.FailRespone(responses.NewInternalServerError()))
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse(gin.H{"is_read": true}, http.StatusOK))
}
