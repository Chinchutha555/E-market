package handlers

import (
	"emarket-api/middlewares"
	"emarket-api/models"
	"emarket-api/responses"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strconv"
)

func SignedDocumentUpload(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		var body struct{ File string }
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 15*1024*1024)
		if c.ShouldBindJSON(&body) != nil || len(body.File) > 14*1024*1024 {
			c.JSON(422, responses.FailRespone(responses.NewAppErr(422, "ไฟล์ไม่ถูกต้อง")))
			return
		}
		data, err := base64.StdEncoding.DecodeString(body.File)
		if err != nil || len(data) == 0 || len(data) > 10*1024*1024 {
			c.JSON(422, responses.FailRespone(responses.NewAppErr(422, "ไฟล์ต้องมีขนาดไม่เกิน 10 MB")))
			return
		}
		kind := http.DetectContentType(data)
		if kind != "application/pdf" && kind != "image/jpeg" && kind != "image/png" {
			c.JSON(422, responses.FailRespone(responses.NewAppErr(422, "รองรับไฟล์ PDF, JPG หรือ PNG")))
			return
		}
		conflict := errors.New("application not awaiting document")
		err = db.Transaction(func(tx *gorm.DB) error {
			var app models.Application
			if err := tx.Where("id = ? AND status = ? AND application_state = ?", id, "A", "approved").First(&app).Error; err != nil {
				return conflict
			}
			result := tx.Model(&models.Application{}).Where("id = ? AND status = ? AND application_state = ?", id, "A", "approved").Updates(map[string]interface{}{"signed_document": data, "signed_document_type": kind, "application_state": "completed", "updated_by": middlewares.GetUserID(c)})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return conflict
			}
			if app.CreatedBy != nil && app.ApplicationType != "general" {
				n := models.Notification{RecipientRole: "tenant", RecipientUserID: app.CreatedBy, Type: "contract_ready", StatusCode: "completed", StatusLabel: "พร้อมทำสัญญา", Title: "คำร้องเสร็จสมบูรณ์", Message: "เอกสารคำร้องครบแล้ว กรุณาจัดทำสัญญาเช่า", TargetURL: fmt.Sprintf("/contracts?new_for_stall=%s&app_id=%d", url.QueryEscape(app.StallIdent), app.ID)}
				if app.BookingID != nil {
					var booking models.Booking
					lookupErr := tx.Select("id", "requires_contract").First(&booking, *app.BookingID).Error
					if lookupErr != nil && !errors.Is(lookupErr, gorm.ErrRecordNotFound) {
						return lookupErr
					}
					if lookupErr == nil && booking.RequiresContract != nil && !*booking.RequiresContract {
						n.Type = "application_completed"
						n.StatusLabel = "ดำเนินการเสร็จแล้ว"
						n.Title = "คำร้องเสร็จสมบูรณ์"
						n.Message = "เอกสารคำร้องครบแล้ว ไม่ต้องทำสัญญา"
						n.TargetURL = "/applications"
					}
				}
				if err := tx.Create(&n).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if errors.Is(err, conflict) {
			c.JSON(409, responses.FailRespone(responses.NewAppErr(409, "คำร้องไม่อยู่ในขั้นตอนรับเอกสาร")))
			return
		}
		if err != nil {
			c.JSON(500, responses.FailRespone(responses.NewInternalServerError()))
			return
		}
		c.JSON(200, responses.SuccessResponse(gin.H{"status": "completed"}, 200))
	}
}

func SignedDocumentDownload(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatus(400)
			return
		}
		var app models.Application
		if db.Select("signed_document", "signed_document_type").Where("id = ? AND status = ?", id, "A").First(&app).Error != nil || len(app.SignedDocument) == 0 {
			c.AbortWithStatus(404)
			return
		}
		c.Header("Content-Disposition", "attachment; filename=\"signed-application."+map[string]string{"application/pdf": "pdf", "image/jpeg": "jpg", "image/png": "png"}[app.SignedDocumentType]+"\"")
		c.Data(200, app.SignedDocumentType, app.SignedDocument)
	}
}
