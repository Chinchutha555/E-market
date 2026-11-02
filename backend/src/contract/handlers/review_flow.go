package handlers

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"emarket-api/middlewares"
	"emarket-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var errContractForbidden = errors.New("forbidden")
var errContractInvalid = errors.New("invalid")
var errContractConflict = errors.New("conflict")

func contractFile(value string, imageOnly bool) ([]byte, string, bool) {
	if len(value) == 0 || len(value) > 4*1024*1024 {
		return nil, "", false
	}
	b, err := base64.StdEncoding.DecodeString(value)
	if err != nil || len(b) == 0 || len(b) > 3*1024*1024 {
		return nil, "", false
	}
	kind := http.DetectContentType(b)
	return b, kind, kind == "image/png" || kind == "image/jpeg" || (!imageOnly && kind == "application/pdf")
}

func notifyContract(tx *gorm.DB, contract *models.Contract, role, title, message string) error {
	var users []models.User
	query := tx.Where("role = ? AND status = ?", role, "A")
	if role == "staff" || role == "market_director" {
		if contract.MarketID == nil {
			return nil
		}
		query = query.Where("id IN (SELECT user_id FROM market_staff_assignment WHERE market_id = ?)", *contract.MarketID)
	}
	if role == "tenant" {
		if contract.TenantUserID == nil {
			return nil
		}
		query = query.Where("id = ?", *contract.TenantUserID)
	}
	if err := query.Find(&users).Error; err != nil {
		return err
	}
	for _, user := range users {
		id := user.ID
		recipientRole := role
		if role == "staff" {
			recipientRole = "officer"
		}
		n := models.Notification{RecipientRole: recipientRole, RecipientUserID: &id, Type: "contract_workflow", StatusCode: contract.ContractState, StatusLabel: title, Title: title, Message: message, StallNumber: contract.StallIdent, TargetURL: "/contracts"}
		if err := tx.Create(&n).Error; err != nil {
			return err
		}
	}
	return nil
}

func ContractReview(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		var req struct {
			Action    string
			Signature string
			OTP       string
			Reason    string
		}
		if err != nil || c.ShouldBindJSON(&req) != nil {
			c.JSON(400, gin.H{"message": "ข้อมูลไม่ถูกต้อง"})
			return
		}
		role, userID := middlewares.GetUserRole(c), middlewares.GetUserID(c)
		err = db.Transaction(func(tx *gorm.DB) error {
			var contract models.Contract
			if err := tx.Where("id = ? AND status = ?", id, "A").First(&contract).Error; err != nil {
				return err
			}
			allowed := (contract.ContractState == "pending" && role == "staff") || (contract.ContractState == "pending_market_director" && role == "market_director") || (contract.ContractState == "pending_bureau_director" && role == "bureau_director")
			if !allowed {
				return errContractForbidden
			}
			if req.Action == "approve" {
				if len(contract.TenantSignature) == 0 || len(contract.SupportingDocument) == 0 {
					return errContractInvalid
				}
				if role == "market_director" && len(contract.OfficerSignature) == 0 {
					return errContractInvalid
				}
				if role == "bureau_director" && (len(contract.OfficerSignature) == 0 || len(contract.MarketDirectorSignature) == 0) {
					return errContractInvalid
				}
			}
			updates := map[string]interface{}{"updated_by": userID}
			switch req.Action {
			case "return_revision":
				if role == "bureau_director" || strings.TrimSpace(req.Reason) == "" {
					return errContractInvalid
				}
				updates["contract_state"], updates["review_reason"] = "revision_required", strings.TrimSpace(req.Reason)
			case "approve":
				signature, kind, ok := contractFile(req.Signature, true)
				if !ok {
					return errContractInvalid
				}
				switch role {
				case "staff":
					updates["officer_signature"], updates["officer_signature_type"], updates["contract_state"] = signature, kind, "pending_market_director"
				case "market_director":
					updates["market_director_signature"], updates["market_director_signature_type"], updates["contract_state"] = signature, kind, "pending_bureau_director"
				case "bureau_director":
					// Prototype OTP; replace with a real email challenge before production.
					if req.OTP != "123456" {
						return errContractInvalid
					}
					updates["bureau_director_signature"], updates["bureau_director_signature_type"], updates["contract_state"] = signature, kind, "active"
				}
			default:
				return errContractInvalid
			}
			res := tx.Model(&models.Contract{}).Where("id = ? AND contract_state = ? AND status = ?", id, contract.ContractState, "A").Updates(updates)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected != 1 {
				return errContractConflict
			}
			contract.ContractState = updates["contract_state"].(string)
			switch contract.ContractState {
			case "pending_market_director":
				return notifyContract(tx, &contract, "market_director", "รอตรวจสัญญา", "เจ้าหน้าที่ตรวจสัญญาแล้ว กรุณาพิจารณาและลงนาม")
			case "pending_bureau_director":
				return notifyContract(tx, &contract, "bureau_director", "รออนุมัติสัญญา", "ผอ.ตลาดลงนามแล้ว กรุณาพิจารณาและลงนาม")
			case "revision_required":
				return notifyContract(tx, &contract, "tenant", "สัญญาส่งกลับแก้ไข", strings.TrimSpace(req.Reason))
			case "active":
				var existing int64
				if err := tx.Model(&models.Contract{}).Where("stall_ident = ? AND contract_state = ? AND status = ? AND id <> ?", contract.StallIdent, "active", "A", contract.ID).Count(&existing).Error; err != nil {
					return err
				}
				if existing > 0 {
					return errContractConflict
				}
				if err := tx.Model(&models.Stall{}).Where("ident = ?", contract.StallIdent).Update("stall_status", "occupied").Error; err != nil {
					return err
				}
				if err := tx.Model(&models.Booking{}).Where("stall_ident = ? AND status = ? AND booking_status = ?", contract.StallIdent, "A", "lease_applied").Update("booking_status", "completed").Error; err != nil {
					return err
				}
				var waiters []models.Waitlist
				if err := tx.Where("stall_ident = ? AND status IN ?", contract.StallIdent, []string{"waiting", "offered"}).Find(&waiters).Error; err != nil {
					return err
				}
				for _, waiter := range waiters {
					if err := tx.Model(&models.Waitlist{}).Where("id = ?", waiter.ID).Updates(map[string]interface{}{"status": "cancelled", "dismiss_reason": "แผงมีสัญญาที่อนุมัติแล้ว", "expires_at": nil}).Error; err != nil {
						return err
					}
					id := waiter.TenantID
					notification := models.Notification{RecipientRole: "tenant", RecipientUserID: &id, Type: "waitlist_cancelled", StatusCode: "waitlist_cancelled", Title: "แผงมีผู้เช่าแล้ว", Message: "แผงที่รอมีสัญญาอนุมัติแล้ว กรุณาเลือกแผงอื่น", StallNumber: contract.StallIdent, TargetURL: "/bookings?tab=interests"}
					if err := tx.Create(&notification).Error; err != nil {
						return err
					}
				}
				return notifyContract(tx, &contract, "tenant", "สัญญาอนุมัติแล้ว", "สัญญาเช่าอนุมัติครบแล้ว กรุณาดำเนินการชำระเงิน")
			}
			return nil
		})
		switch err {
		case nil:
			c.JSON(200, gin.H{"success": true})
		case errContractForbidden:
			c.JSON(403, gin.H{"message": "ไม่มีสิทธิ์ดำเนินการในขั้นตอนนี้"})
		case errContractInvalid:
			c.JSON(422, gin.H{"message": "กรุณาตรวจสอบลายเซ็น รหัส OTP หรือเหตุผล"})
		case errContractConflict:
			c.JSON(409, gin.H{"message": "สถานะสัญญาเปลี่ยนแล้ว"})
		default:
			c.JSON(500, gin.H{"message": "บันทึกสัญญาไม่สำเร็จ"})
		}
	}
}

func ContractResubmit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		var req struct {
			TenantName, TenantPhone, TenantIDCard, TenantEmail, TenantAddress, ShopName, BusinessType, BillingCycle, TermsAndConditions, StartDate, EndDate, Signature, Document string
			MonthlyRental, SecurityDeposit, AdvancePayment                                                                                                                       float64
		}
		if err != nil || c.ShouldBindJSON(&req) != nil {
			c.JSON(400, gin.H{"message": "ข้อมูลไม่ถูกต้อง"})
			return
		}
		sig, sigType, ok := contractFile(req.Signature, true)
		doc, docType, docOK := contractFile(req.Document, false)
		if !ok || !docOK || strings.TrimSpace(req.TenantName) == "" || strings.TrimSpace(req.TenantPhone) == "" || strings.TrimSpace(req.ShopName) == "" || req.StartDate == "" || req.EndDate == "" || req.MonthlyRental <= 0 {
			c.JSON(422, gin.H{"message": "กรุณากรอกข้อมูล แนบเอกสารและลงลายเซ็น"})
			return
		}
		start, startErr := time.Parse("2006-01-02", req.StartDate)
		end, endErr := time.Parse("2006-01-02", req.EndDate)
		if startErr != nil || endErr != nil || !end.After(start) || req.SecurityDeposit < 0 || req.AdvancePayment < 0 || strings.TrimSpace(req.TenantIDCard) == "" {
			c.JSON(422, gin.H{"message": "กรุณาตรวจสอบวันที่และยอดเงิน"})
			return
		}
		err = db.Transaction(func(tx *gorm.DB) error {
			var contract models.Contract
			if err := tx.Where("id = ? AND tenant_user_id = ? AND contract_state = ? AND status = ?", id, middlewares.GetUserID(c), "revision_required", "A").First(&contract).Error; err != nil {
				return errContractConflict
			}
			updates := map[string]interface{}{"tenant_name": strings.TrimSpace(req.TenantName), "tenant_phone": strings.TrimSpace(req.TenantPhone), "tenant_id_card": strings.TrimSpace(req.TenantIDCard), "tenant_email": strings.TrimSpace(req.TenantEmail), "tenant_address": strings.TrimSpace(req.TenantAddress), "business_type": strings.TrimSpace(req.BusinessType), "billing_cycle": req.BillingCycle, "terms_and_conditions": req.TermsAndConditions, "advance_payment": req.AdvancePayment, "shop_name": strings.TrimSpace(req.ShopName), "start_date": req.StartDate, "end_date": req.EndDate, "monthly_rental": req.MonthlyRental, "security_deposit": req.SecurityDeposit, "tenant_signature": sig, "tenant_signature_type": sigType, "supporting_document": doc, "supporting_document_type": docType, "officer_signature": nil, "officer_signature_type": "", "market_director_signature": nil, "market_director_signature_type": "", "bureau_director_signature": nil, "bureau_director_signature_type": "", "review_reason": "", "contract_state": "pending", "updated_by": middlewares.GetUserID(c)}
			res := tx.Model(&models.Contract{}).Where("id = ? AND contract_state = ?", id, "revision_required").Updates(updates)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected != 1 {
				return errContractConflict
			}
			contract.ContractState = "pending"
			return notifyContract(tx, &contract, "staff", "รอตรวจสัญญา", "ผู้เช่าแก้ไขและส่งสัญญาอีกครั้งแล้ว")
		})
		if err != nil {
			c.JSON(409, gin.H{"message": "ไม่สามารถส่งสัญญาอีกครั้งได้"})
			return
		}
		c.JSON(200, gin.H{"success": true})
	}
}

func ContractAsset(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		var contract models.Contract
		if err != nil || db.Where("id = ? AND status = ?", id, "A").First(&contract).Error != nil {
			c.AbortWithStatus(404)
			return
		}
		if contract.ContractState == "draft" && middlewares.GetUserRole(c) != "staff" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		var data []byte
		var kind string
		switch c.Param("kind") {
		case "tenant":
			data, kind = contract.TenantSignature, contract.TenantSignatureType
		case "officer":
			data, kind = contract.OfficerSignature, contract.OfficerSignatureType
		case "market":
			data, kind = contract.MarketDirectorSignature, contract.MarketDirectorSignatureType
		case "bureau":
			data, kind = contract.BureauDirectorSignature, contract.BureauDirectorSignatureType
		case "document":
			data, kind = contract.SupportingDocument, contract.SupportingDocumentType
		default:
			c.AbortWithStatus(404)
			return
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
