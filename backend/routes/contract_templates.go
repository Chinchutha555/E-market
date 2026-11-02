package routes

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"emarket-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var contractTemplateCode = regexp.MustCompile("^[a-z0-9_-]{1,40}$")

func contractTemplates(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		marketID, err := strconv.Atoi(c.Param("id"))
		if err != nil || marketID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "ตลาดไม่ถูกต้อง"})
			return
		}
		if c.Request.Method == http.MethodGet {
			rows := []models.MarketContractTemplate{}
			if err := db.Where("market_id = ?", marketID).Order("sort_order ASC, id ASC").Find(&rows).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "ไม่สามารถโหลดรูปแบบสัญญาได้"})
				return
			}
			if marketID == 1 {
				hasJJ := false
				for _, r := range rows {
					if r.Code == "jj_semi_permanent_lease" {
						hasJJ = true
						break
					}
				}
				if !hasJJ {
					jjTemplate := models.MarketContractTemplate{
						MarketID:     1,
						Code:         "jj_semi_permanent_lease",
						Title:        "สัญญาให้เช่าแผงค้าภายในบริเวณตลาดนัดจตุจักร (แผงค้ากึ่งถาวร)",
						ContractKind: "standard",
						Description:  "บันทึกข้อตกลงสัญญาให้เช่าแผงค้ากึ่งถาวรภายในบริเวณตลาดนัดจตุจักร สำนักงานตลาดกรุงเทพมหานคร กำหนดระยะเวลา 1 ปี",
						IsActive:     true,
						SortOrder:    1,
					}
					if err := db.Create(&jjTemplate).Error; err == nil {
						rows = append([]models.MarketContractTemplate{jjTemplate}, rows...)
					}
				}
			}
			c.JSON(http.StatusOK, gin.H{"data": rows})
			return
		}

		var body struct {
			Templates []models.MarketContractTemplate `json:"templates"`
		}
		if c.ShouldBindJSON(&body) != nil || body.Templates == nil || len(body.Templates) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "ข้อมูลรูปแบบสัญญาไม่ถูกต้อง"})
			return
		}
		kinds := map[string]bool{"standard": true, "trial": true, "temporary": true, "custom": true}
		seenCodes, seenTitles := map[string]bool{}, map[string]bool{}
		for i := range body.Templates {
			row := &body.Templates[i]
			row.MarketID = marketID
			row.Code = strings.ToLower(strings.TrimSpace(row.Code))
			row.Title = strings.TrimSpace(row.Title)
			row.ContractKind = strings.ToLower(strings.TrimSpace(row.ContractKind))
			row.Description = strings.TrimSpace(row.Description)
			titleKey := strings.ToLower(row.Title)
			if !contractTemplateCode.MatchString(row.Code) || row.Title == "" || len([]rune(row.Title)) > 255 ||
				!kinds[row.ContractKind] || seenCodes[row.Code] || seenTitles[titleKey] {
				c.JSON(http.StatusBadRequest, gin.H{"message": "กรุณาตรวจสอบชื่อ รหัส และประเภทสัญญา"})
				return
			}
			if len([]rune(row.Description)) > 1000 || row.SortOrder < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"message": "รายละเอียดรูปแบบสัญญาไม่ถูกต้อง"})
				return
			}
			row.ID = 0
			seenCodes[row.Code], seenTitles[titleKey] = true, true
		}
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("market_id = ?", marketID).Delete(&models.MarketContractTemplate{}).Error; err != nil {
				return err
			}
			if len(body.Templates) > 0 {
				return tx.Create(&body.Templates).Error
			}
			return nil
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "ไม่สามารถบันทึกรูปแบบสัญญาได้"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": body.Templates})
	}
}
