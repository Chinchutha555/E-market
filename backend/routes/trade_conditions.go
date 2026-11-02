package routes

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"emarket-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type tradeCondition = models.MarketTradeCondition

var tradeTypeCode = regexp.MustCompile("^[a-z0-9_-]{1,20}$")

func tradeConditions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid market"})
			return
		}
		if c.Request.Method == http.MethodGet {
			rows := []tradeCondition{}
			if err := db.Where("market_id = ?", id).Order("trade_type").Find(&rows).Error; err != nil {
				c.JSON(500, gin.H{"message": "Cannot load settings"})
				return
			}
			for i := range rows {
				rows[i].Options = rows[i].DurationOptions()
			}
			c.JSON(200, gin.H{"data": rows})
			return
		}
		var body struct{ Conditions []tradeCondition }
		if c.ShouldBindJSON(&body) != nil || body.Conditions == nil || len(body.Conditions) > 100 {
			c.JSON(400, gin.H{"message": "Invalid settings"})
			return
		}
		seenCodes, seenTitles := map[string]bool{}, map[string]bool{}
		for i := range body.Conditions {
			row := &body.Conditions[i]
			if row.RequiresContract == nil {
				yes := true
				row.RequiresContract = &yes
			}
			row.Title = strings.TrimSpace(row.Title)
			titleKey := strings.ToLower(row.Title)
			if !tradeTypeCode.MatchString(row.TradeType) || seenCodes[row.TradeType] || row.Title == "" || len([]rune(row.Title)) > 255 || seenTitles[titleKey] {
				c.JSON(400, gin.H{"message": "Invalid settings"})
				return
			}
			options := row.Options
			if options == nil {
				options = row.DurationOptions()
			}
			if len(options) == 0 || len(options) > 100 {
				c.JSON(400, gin.H{"message": "Invalid duration options"})
				return
			}
			seenOptions := map[string]bool{}
			for _, option := range options {
				if !tradeTypeCode.MatchString(option.ID) || seenOptions[option.ID] || option.Duration < 1 || option.Duration > 10000 || (option.Unit != "days" && option.Unit != "months" && option.Unit != "years") {
					c.JSON(400, gin.H{"message": "Invalid duration option"})
					return
				}
				if option.StartDate != "" || option.EndDate != "" {
					start, a := time.Parse("2006-01-02", option.StartDate)
					end, b := time.Parse("2006-01-02", option.EndDate)
					if a != nil || b != nil || end.Before(start) {
						c.JSON(400, gin.H{"message": "Invalid contract dates"})
						return
					}
				}
				seenOptions[option.ID] = true
			}
			encoded, err := json.Marshal(options)
			if err != nil {
				c.JSON(400, gin.H{"message": "Invalid duration options"})
				return
			}
			row.Options, row.OptionsJSON = options, string(encoded)
			row.MinDuration, row.MaxDuration, row.Unit = 1, options[0].Duration, options[0].Unit
			row.StartDate, row.EndDate = options[0].StartDate, options[0].EndDate
			row.MarketID = id
			seenCodes[row.TradeType], seenTitles[titleKey] = true, true
		}
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("market_id = ?", id).Delete(&tradeCondition{}).Error; err != nil {
				return err
			}
			if len(body.Conditions) > 0 {
				return tx.Create(&body.Conditions).Error
			}
			return nil
		}); err != nil {
			c.JSON(500, gin.H{"message": "Cannot save settings"})
			return
		}
		c.JSON(200, gin.H{"data": body.Conditions})
	}
}
