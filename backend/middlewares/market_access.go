package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"emarket-api/models"
	"emarket-api/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RequireBodyMarketAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserRole(c) == "admin" {
			c.Next()
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			forbidMarket(c)
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		var payload struct {
			MarketID int `json:"market_id"`
		}
		if json.Unmarshal(body, &payload) != nil || !CanManageMarket(db, GetUserID(c), GetUserRole(c), payload.MarketID) {
			forbidMarket(c)
			return
		}
		c.Next()
	}
}

// ScopeMarketQuery forces list endpoints into the authenticated account's
// server-side scope. Any market_id/created_by sent by the browser is replaced.
func ScopeMarketQuery(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		switch GetUserRole(c) {
		case "staff":
			ids, err := StaffMarketIDs(db, GetUserID(c))
			if err != nil || len(ids) == 0 {
				query.Set("market_id", "-1")
			} else {
				requestedMarket := query.Get("market_id")
				isAssigned := false
				if requestedMarket != "" {
					for _, id := range ids {
						if strconv.Itoa(id) == requestedMarket {
							isAssigned = true
							break
						}
					}
				}
				if !isAssigned {
					query.Set("market_id", strconv.Itoa(ids[0]))
				}
			}
			query.Del("created_by")
		case "tenant":
			query.Set("created_by", strconv.Itoa(GetUserID(c)))
			query.Del("market_id")
		default:
			query.Del("created_by")
		}
		c.Request.URL.RawQuery = query.Encode()
		c.Next()
	}
}

type scopedResource struct {
	MarketID *int
	StallID  *int
	Ident    string
	OwnerID  *int
}

func resolveResource(db *gorm.DB, kind string, id int) (*scopedResource, error) {
	resource := &scopedResource{}
	switch kind {
	case "booking":
		var item models.Booking
		if err := db.Select("market_id", "stall_id", "stall_ident", "created_by").First(&item, id).Error; err != nil {
			return nil, err
		}
		resource.MarketID, resource.StallID, resource.Ident, resource.OwnerID = item.MarketID, item.StallID, item.StallIdent, item.CreatedBy
	case "application":
		var item models.Application
		if err := db.Select("market_id", "stall_id", "stall_ident", "created_by").First(&item, id).Error; err != nil {
			return nil, err
		}
		resource.MarketID, resource.StallID, resource.Ident, resource.OwnerID = item.MarketID, item.StallID, item.StallIdent, item.CreatedBy
	case "contract":
		var item models.Contract
		if err := db.Select("market_id", "stall_id", "stall_ident", "tenant_user_id").First(&item, id).Error; err != nil {
			return nil, err
		}
		resource.MarketID, resource.StallID, resource.Ident, resource.OwnerID = item.MarketID, item.StallID, item.StallIdent, item.TenantUserID
	default:
		return nil, gorm.ErrRecordNotFound
	}
	if resource.MarketID == nil || *resource.MarketID <= 0 {
		var stall models.Stall
		query := db.Select("market_id").Where("status = 'A'")
		if resource.StallID != nil && *resource.StallID > 0 {
			query = query.Where("id = ?", *resource.StallID)
		} else {
			query = query.Where("ident = ?", resource.Ident)
		}
		if query.First(&stall).Error == nil {
			resource.MarketID = &stall.MarketID
		}
	}
	return resource, nil
}

func RequireResourceAccess(db *gorm.DB, kind string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserRole(c) == "admin" {
			c.Next()
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			forbidMarket(c)
			return
		}
		resource, err := resolveResource(db, kind, id)
		if err != nil {
			forbidMarket(c)
			return
		}
		if GetUserRole(c) == "tenant" {
			if resource.OwnerID == nil || *resource.OwnerID != GetUserID(c) {
				forbidMarket(c)
				return
			}
			c.Next()
			return
		}
		if resource.MarketID == nil || !CanManageMarket(db, GetUserID(c), GetUserRole(c), *resource.MarketID) {
			forbidMarket(c)
			return
		}
		c.Next()
	}
}

func RequireStallPayloadAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserRole(c) == "admin" {
			c.Next()
			return
		}
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			forbidMarket(c)
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		var payload struct {
			StallID     int    `json:"stall_id"`
			StallIdent  string `json:"stall_ident"`
			StallNumber string `json:"stall_number"`
		}
		if json.Unmarshal(body, &payload) != nil {
			forbidMarket(c)
			return
		}
		ident := payload.StallIdent
		if ident == "" {
			ident = payload.StallNumber
		}
		var stall models.Stall
		query := db.Select("market_id").Where("status = 'A'")
		if payload.StallID > 0 {
			query = query.Where("id = ?", payload.StallID)
		} else {
			query = query.Where("ident = ?", ident)
		}
		if query.First(&stall).Error != nil || !CanManageMarket(db, GetUserID(c), GetUserRole(c), stall.MarketID) {
			forbidMarket(c)
			return
		}
		c.Next()
	}
}

func StaffMarketIDs(db *gorm.DB, userID int) ([]int, error) {
	var ids []int
	err := db.Model(&models.MarketStaffAssignment{}).
		Where("user_id = ?", userID).Order("market_id ASC").Pluck("market_id", &ids).Error
	return ids, err
}

func CanManageMarket(db *gorm.DB, userID int, role string, marketID int) bool {
	if role == "admin" {
		return true
	}
	if role != "staff" || userID <= 0 || marketID <= 0 {
		return false
	}
	var count int64
	db.Model(&models.MarketStaffAssignment{}).
		Where("user_id = ? AND market_id = ?", userID, marketID).Count(&count)
	return count > 0
}

func forbidMarket(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, responses.FailRespone(
		responses.NewAppErr(http.StatusForbidden, "คุณไม่มีสิทธิ์จัดการตลาดนี้"),
	))
}

// RequireMarketAccess protects market-scoped management endpoints. paramName
// may contain either a numeric market ID or a market slug.
func RequireMarketAccess(db *gorm.DB, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserRole(c) == "admin" {
			c.Next()
			return
		}
		value := c.Param(paramName)
		var market models.Market
		query := db.Where("status != 'D'")
		if id, err := strconv.Atoi(value); err == nil {
			query = query.Where("id = ?", id)
		} else {
			query = query.Where("slug = ?", value)
		}
		if query.First(&market).Error != nil || !CanManageMarket(db, GetUserID(c), GetUserRole(c), market.ID) {
			forbidMarket(c)
			return
		}
		c.Set("marketID", market.ID)
		c.Next()
	}
}

func RequireStallAccess(db *gorm.DB, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetUserRole(c) == "admin" {
			c.Next()
			return
		}
		stallID, err := strconv.Atoi(c.Param(paramName))
		if err != nil {
			forbidMarket(c)
			return
		}
		var stall models.Stall
		if db.Select("id", "market_id").First(&stall, stallID).Error != nil ||
			!CanManageMarket(db, GetUserID(c), GetUserRole(c), stall.MarketID) {
			forbidMarket(c)
			return
		}
		c.Next()
	}
}
