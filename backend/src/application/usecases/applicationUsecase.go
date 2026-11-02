package usecases

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"emarket-api/logs"
	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/application/domains"
	"gorm.io/gorm"
)

func normalizedApplicationState(s string) string {
	if s == "" || s == "submitted" {
		return "pending_officer"
	}
	return s
}

type useCase struct {
	repo domains.Repository
}

func NewUseCase(repo domains.Repository) domains.UseCase {
	return &useCase{repo: repo}
}

func toApplicationResponse(a *models.Application) responses.ApplicationResponse {
	monthlyRate := a.MonthlyRate
	if monthlyRate <= 0 && a.ApplicationType != "general" && !strings.HasPrefix(a.ApplicationNo, "GEN-") {
		monthlyRate = 3600
	}
	periodMonths := a.LeaseTermMonths
	if periodMonths <= 0 && a.TradeDurationUnit != "days" {
		periodMonths = 12
	}
	startDate := a.RequestedStartDate
	if startDate == "" && a.TradeOptionID == "" {
		startDate = a.CreatedAt.AddDate(0, 0, 7).Format("2006-01-02")
	}

	appType := a.ApplicationType
	if appType == "" {
		if strings.HasPrefix(a.ApplicationNo, "GEN-") {
			appType = "general"
		} else if a.BookingID != nil || (a.StallIdent != "" && a.StallIdent != "-") {
			appType = "stall_lease"
		} else {
			appType = "general"
		}
	}
	subject := a.Subject
	if subject == "" && appType == "general" {
		subject = a.ShopName
	}

	var marketDirectorReview *models.MarketDirectorReview
	if a.MarketDirectorReviewJSON != "" {
		var review models.MarketDirectorReview
		if json.Unmarshal([]byte(a.MarketDirectorReviewJSON), &review) == nil {
			marketDirectorReview = &review
		}
	}
	var bureauDirectorReview *models.BureauDirectorReview
	if a.BureauDirectorReviewJSON != "" {
		var review models.BureauDirectorReview
		if json.Unmarshal([]byte(a.BureauDirectorReviewJSON), &review) == nil {
			bureauDirectorReview = &review
		}
	}
	var officerReview *models.OfficerReview
	if a.OfficerReviewJSON != "" {
		var review models.OfficerReview
		if json.Unmarshal([]byte(a.OfficerReviewJSON), &review) == nil {
			officerReview = &review
		}
	}
	return responses.ApplicationResponse{
		DirectorReviewSupported:    true,
		MarketDirectorReview:       marketDirectorReview,
		BureauDirectorReview:       bureauDirectorReview,
		OfficerReviewSupported:     true,
		OfficerReview:              officerReview,
		TradeOptionID:              a.TradeOptionID,
		TradeStartDate:             a.TradeStartDate,
		TradeEndDate:               a.TradeEndDate,
		TradeDuration:              a.TradeDuration,
		TradeDurationUnit:          a.TradeDurationUnit,
		RequiresContract:           a.RequiresContract,
		ID:                         a.ID,
		ApplicationNo:              a.ApplicationNo,
		ApplicationNumber:          a.ApplicationNo,
		ApplicationType:            appType,
		Subject:                    subject,
		BookingID:                  a.BookingID,
		BookingNumber:              a.BookingNumber,
		DocumentNumber:             a.DocumentNumber,
		TenantSignature:            a.TenantSignature,
		HasOfficerSignature:        len(a.OfficerSignature) > 0,
		HasFinanceSignature:        len(a.FinanceSignature) > 0,
		HasMarketDirectorSignature: len(a.MarketDirectorSignature) > 0,
		HasBureauDirectorSignature: len(a.BureauDirectorSignature) > 0,
		HasSignedDocument:          len(a.SignedDocument) > 0,
		TenantSignedAt:             a.TenantSignedAt,
		StallID:                    a.StallID,
		MarketID:                   a.MarketID,
		MarketName:                 a.MarketName,
		StallIdent:                 a.StallIdent,
		StallNumber:                a.StallIdent,
		ApplicantName:              a.ApplicantName,
		TenantName:                 a.ApplicantName,
		ApplicantPhone:             a.ApplicantPhone,
		TenantPhone:                a.ApplicantPhone,
		ApplicantIDCard:            a.ApplicantIDCard,
		ApplicantAge:               a.ApplicantAge,
		HasIDCardCopy:              len(a.IDCardCopy) > 0,
		ApplicantEmail:             a.ApplicantEmail,
		ApplicantAddress:           a.ApplicantAddress,
		Details:                    a.Details,
		ShopName:                   a.ShopName,
		ProductCategory:            a.ProductCategory,
		ProductType:                a.ProductCategory,
		ProposedRentType:           a.ProposedRentType,
		RequestedStartDate:         startDate,
		LeaseTermMonths:            periodMonths,
		RequestedPeriodMonths:      periodMonths,
		MonthlyRate:                monthlyRate,
		StallMonthlyRate:           monthlyRate,
		Status:                     normalizedApplicationState(a.ApplicationState),
		OfficerComment:             a.OfficerComment,
		MarketDirectorComment:      a.MarketDirectorComment,
		BureauDirectorComment:      a.BureauDirectorComment,
		ReviewedAt:                 a.ReviewedAt,
		ReviewNote:                 a.ReviewNote,
		CreatedAt:                  a.CreatedAt,
	}
}

func (u *useCase) GetApplications(filter requests.ApplicationFilter) ([]responses.ApplicationResponse, error) {
	apps, err := u.repo.FindAll(filter)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	res := make([]responses.ApplicationResponse, len(apps))
	for i, a := range apps {
		res[i] = toApplicationResponse(&a)
	}
	return res, nil
}

func (u *useCase) GetApplicationByID(id int) (*responses.ApplicationResponse, error) {
	app, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลคำร้อง")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	res := toApplicationResponse(app)
	return &res, nil
}

func (u *useCase) CreateApplication(userID int, req requests.CreateApplicationRequest) (*responses.ApplicationResponse, error) {
	if req.ApplicantAge != nil && (*req.ApplicantAge < 1 || *req.ApplicantAge > 120) {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุอายุระหว่าง 1 ถึง 120 ปี")
	}
	if strings.TrimSpace(req.IDCardCopyBase64) == "" {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาแนบสำเนาบัตรประชาชน")
	}
	var idCardCopy []byte
	var idCardCopyType string
	if req.IDCardCopyBase64 != "" {
		if len(req.IDCardCopyBase64) > 4*1024*1024 {
			return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "ไฟล์สำเนาบัตรประชาชนต้องไม่เกิน 3 MB")
		}
		var err error
		idCardCopy, err = base64.StdEncoding.DecodeString(req.IDCardCopyBase64)
		if err != nil || len(idCardCopy) == 0 || len(idCardCopy) > 3*1024*1024 {
			return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "ไฟล์สำเนาบัตรประชาชนไม่ถูกต้อง")
		}
		idCardCopyType = http.DetectContentType(idCardCopy)
		if idCardCopyType != "application/pdf" && idCardCopyType != "image/jpeg" && idCardCopyType != "image/png" {
			return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "รองรับไฟล์ PDF, JPG หรือ PNG เท่านั้น")
		}
	}
	appType := strings.TrimSpace(req.ApplicationType)
	var parsedBookingID *int
	switch v := req.BookingID.(type) {
	case float64:
		idVal := int(v)
		if idVal > 0 {
			parsedBookingID = &idVal
		}
	case int:
		if v > 0 {
			parsedBookingID = &v
		}
	}

	stallIdent := strings.TrimSpace(req.StallIdent)
	if stallIdent == "" {
		stallIdent = strings.TrimSpace(req.StallNumber)
	}

	if appType == "" {
		if parsedBookingID != nil {
			appType = "stall_lease"
		} else if stallIdent != "" && stallIdent != "-" {
			appType = "stall_lease"
		} else {
			appType = "general"
		}
	}

	if appType == "stall_lease" {
		if stallIdent == "" {
			return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุรหัสแผงค้า")
		}
	} else {
		if stallIdent == "" {
			stallIdent = "-"
		}
	}

	applicantName := strings.TrimSpace(req.ApplicantName)
	if applicantName == "" {
		applicantName = strings.TrimSpace(req.TenantName)
	}
	applicantPhone := strings.TrimSpace(req.ApplicantPhone)
	if applicantPhone == "" {
		applicantPhone = strings.TrimSpace(req.TenantPhone)
	}
	if applicantName == "" || applicantPhone == "" {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุชื่อและเบอร์โทรศัพท์ผู้ยื่นคำร้อง")
	}
	idCard := strings.TrimSpace(req.ApplicantIDCard)
	if idCard == "" {
		idCard = strings.TrimSpace(req.IDCard)
	}
	address := strings.TrimSpace(req.ApplicantAddress)
	if address == "" {
		address = strings.TrimSpace(req.Address)
	}
	productCategory := strings.TrimSpace(req.ProductCategory)
	if productCategory == "" {
		productCategory = strings.TrimSpace(req.ProductType)
	}
	if productCategory == "" {
		if appType == "general" {
			productCategory = "คำร้องทั่วไป"
		} else {
			productCategory = "สินค้าทั่วไป"
		}
	}
	subject := strings.TrimSpace(req.Subject)
	shopName := strings.TrimSpace(req.ShopName)
	if shopName == "" {
		if subject != "" {
			shopName = subject
		} else {
			shopName = applicantName
		}
	}
	if subject == "" && appType == "general" {
		subject = shopName
	}
	leaseTerm := req.LeaseTermMonths
	if leaseTerm <= 0 {
		leaseTerm = req.RequestedPeriodMonths
	}
	if leaseTerm <= 0 {
		if appType == "general" {
			leaseTerm = 1
		} else {
			leaseTerm = 12
		}
	}
	monthlyRate := req.MonthlyRate
	if monthlyRate <= 0 {
		monthlyRate = req.Fee
	}
	if monthlyRate <= 0 && appType == "stall_lease" {
		monthlyRate = 3600
	}
	startDate := strings.TrimSpace(req.RequestedStartDate)
	if startDate == "" && req.TradeOptionID == "" {
		startDate = time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	}

	var parsedStallID *int
	switch v := req.StallID.(type) {
	case float64:
		idVal := int(v)
		if idVal > 0 {
			parsedStallID = &idVal
		}
	case int:
		if v > 0 {
			parsedStallID = &v
		}
	}

	var parsedMarketID *int
	switch v := req.MarketID.(type) {
	case float64:
		idVal := int(v)
		if idVal > 0 {
			parsedMarketID = &idVal
		}
	case int:
		if v > 0 {
			parsedMarketID = &v
		}
	}

	prefix := "APP"
	if appType == "general" {
		prefix = "GEN"
	}
	appNo := fmt.Sprintf("%s-%s-%06d", prefix, time.Now().Format("20060102"), time.Now().UnixMicro()%1000000)

	reviewNote := strings.TrimSpace(req.ReviewNote)
	details := strings.TrimSpace(req.Details)
	if details == "" {
		details = reviewNote
	}
	if reviewNote == "" {
		reviewNote = details
	}
	marketName := strings.TrimSpace(req.MarketName)
	proposedRentType := strings.TrimSpace(req.ProposedRentType)
	if proposedRentType == "" {
		proposedRentType = "monthly"
	}

	app := models.Application{
		ApplicationNo:      appNo,
		TradeOptionID:      strings.TrimSpace(req.TradeOptionID),
		ApplicationType:    appType,
		Subject:            subject,
		BookingID:          parsedBookingID,
		StallID:            parsedStallID,
		MarketID:           parsedMarketID,
		MarketName:         marketName,
		StallIdent:         stallIdent,
		ApplicantName:      applicantName,
		ApplicantPhone:     applicantPhone,
		ApplicantIDCard:    idCard,
		ApplicantAge:       req.ApplicantAge,
		IDCardCopy:         idCardCopy,
		IDCardCopyType:     idCardCopyType,
		ApplicantEmail:     strings.TrimSpace(req.ApplicantEmail),
		ApplicantAddress:   address,
		ShopName:           shopName,
		ProductCategory:    productCategory,
		RequestedStartDate: startDate,
		LeaseTermMonths:    leaseTerm,
		MonthlyRate:        monthlyRate,
		ProposedRentType:   proposedRentType,
		Details:            details,
		ReviewNote:         reviewNote,
		ApplicationState:   "pending_officer",
	}
	if userID > 0 {
		app.CreatedBy = &userID
		app.UpdatedBy = &userID
	}

	if err := u.repo.Create(&app); err != nil {
		var appErr *responses.AppErr
		if errors.As(err, &appErr) {
			return nil, appErr
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	res := toApplicationResponse(&app)
	return &res, nil
}

func (u *useCase) UpdateStatus(id int, userID int, req requests.UpdateApplicationStatusRequest) (*responses.ApplicationResponse, error) {
	app, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลคำร้อง")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	now := time.Now()
	app.ApplicationState = req.Status
	if req.MarketDirectorReviewJSON != nil {
		app.MarketDirectorReviewJSON = *req.MarketDirectorReviewJSON
	}
	if req.BureauDirectorReviewJSON != nil {
		app.BureauDirectorReviewJSON = *req.BureauDirectorReviewJSON
	}
	if req.OfficerReviewJSON != nil {
		app.OfficerReviewJSON = *req.OfficerReviewJSON
	}
	if len(req.OfficerSignature) > 0 {
		app.OfficerSignature = req.OfficerSignature
		app.OfficerSignatureType = req.OfficerSignatureType
	}
	if len(req.FinanceSignature) > 0 {
		app.FinanceSignature = req.FinanceSignature
		app.FinanceSignatureType = req.FinanceSignatureType
	}
	if len(req.MarketDirectorSignature) > 0 {
		app.MarketDirectorSignature = req.MarketDirectorSignature
		app.MarketDirectorSignatureType = req.MarketDirectorSignatureType
	}
	if len(req.BureauDirectorSignature) > 0 {
		app.BureauDirectorSignature = req.BureauDirectorSignature
		app.BureauDirectorSignatureType = req.BureauDirectorSignatureType
	}
	if req.Reason != "" {
		app.ReviewNote = req.Reason
	}
	if req.OfficerComment != "" {
		app.OfficerComment = req.OfficerComment
	}
	if req.MarketDirectorComment != "" {
		app.MarketDirectorComment = req.MarketDirectorComment
	}
	if req.BureauDirectorComment != "" {
		app.BureauDirectorComment = req.BureauDirectorComment
	}
	app.ReviewedAt = &now
	app.UpdatedBy = &userID

	if err := u.repo.Update(app); err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	if app.StallIdent != "" {
		if req.Status == "approved" {
			// คำร้องที่อนุมัติแล้วเป็นเพียงการกันแผงไว้เพื่อจัดทำสัญญา
			// แผงจะเป็น occupied เมื่อสัญญามีผลเท่านั้น
			_ = u.repo.UpdateStallStatus(app.StallIdent, "reserved")
		} else if req.Status == "rejected" || req.Status == "cancelled" {
			_ = u.repo.UpdateStallStatus(app.StallIdent, "available")
		}
	}

	res := toApplicationResponse(app)
	return &res, nil
}

func (u *useCase) Action(id int, userID int, req requests.ApplicationActionRequest) (*responses.ApplicationResponse, error) {
	app, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลคำร้อง")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	comment := strings.TrimSpace(req.Comment)
	if comment == "" {
		comment = strings.TrimSpace(req.Reason)
	}

	if app.ApplicationState == "pending_officer" || app.ApplicationState == "submitted" {
		action := strings.ToLower(strings.TrimSpace(req.Action))
		if (action == "return_revision" || action == "reject") && comment == "" {
			return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุหมายเหตุ")
		}
		status := ""
		switch action {
		case "approve":
			status = "pending_market_director"
		case "return_revision":
			status = "revision_required"
		case "reject":
			status = "rejected"
		default:
			return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "การดำเนินการไม่ถูกต้อง")
		}
		var signature []byte
		var signatureType string
		var officerReviewJSON *string
		if action == "approve" && (app.ApplicationType == "stall_lease" || (app.ApplicationType == "" && app.StallIdent != "" && app.StallIdent != "-")) {
			var reviewErr error
			officerReviewJSON, reviewErr = encodeOfficerReview(req.OfficerReview)
			if reviewErr != nil {
				return nil, reviewErr
			}
			if req.OfficerReview.OtherChecked || req.OfficerReview.StallStatus == "other" {
				status = "cancelled"
				if comment == "" {
					comment = "ยกเลิกคำร้องจากผลตรวจสอบของเจ้าหน้าที่"
				}
			} else if req.OfficerReview.StallStatus == "has" {
				status = "pending_debt_check"
			}
		}
		if action == "approve" {
			var signatureErr error
			signature, signatureType, signatureErr = decodeApprovalSignature(req.SignatureBase64)
			if signatureErr != nil {
				return nil, signatureErr
			}
		}
		return u.UpdateStatus(id, userID, requests.UpdateApplicationStatusRequest{Status: status, Reason: comment, OfficerComment: comment, OfficerSignature: signature, OfficerSignatureType: signatureType, OfficerReviewJSON: officerReviewJSON})
	}
	if app.ApplicationState == "pending_debt_check" || app.ApplicationState == "debt_pending" {
		if strings.ToLower(strings.TrimSpace(req.Action)) != "approve" {
			return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "ฝ่ายการคลังต้องบันทึกผลตรวจสอบหนี้")
		}
		review, reviewErr := encodeMarketDirectorReview(req.MarketDirectorReview)
		if reviewErr != nil {
			return nil, reviewErr
		}
		signature, signatureType, signatureErr := decodeApprovalSignature(req.SignatureBase64)
		if signatureErr != nil {
			return nil, signatureErr
		}
		nextStatus := "pending_market_director"
		reason := "ฝ่ายการคลังตรวจสอบแล้ว ไม่มีหนี้ค้างชำระ ส่งต่อ ผอ.ตลาด"
		if req.MarketDirectorReview.DebtStatus == "has" {
			nextStatus = "debt_pending"
			reason = fmt.Sprintf("ฝ่ายการคลังตรวจพบหนี้ค้างชำระ จำนวน %.2f บาท กรุณาชำระหนี้ค้างก่อนส่งให้เจ้าหน้าที่ตรวจสอบอีกครั้ง", req.MarketDirectorReview.DebtAmount)
		}
		return u.UpdateStatus(id, userID, requests.UpdateApplicationStatusRequest{
			Status: nextStatus, Reason: reason, MarketDirectorReviewJSON: review,
			FinanceSignature: signature, FinanceSignatureType: signatureType,
		})
	}
	if app.ApplicationState != "pending_market_director" && app.ApplicationState != "pending_bureau_director" {
		return nil, responses.NewAppErr(http.StatusConflict, "คำร้องไม่ได้อยู่ในขั้นที่สามารถพิจารณาได้")
	}
	action := strings.ToLower(strings.TrimSpace(req.Action))
	isBureau := app.ApplicationState == "pending_bureau_director"
	isStallReview := app.ApplicationType == "stall_lease" || (app.ApplicationType == "" && app.StallIdent != "" && app.StallIdent != "-")
	if isBureau {
		if action != "approve" && action != "reject" {
			return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "ผอ.สำนักเลือกได้เฉพาะอนุมัติหรือไม่อนุมัติ")
		}
		// Temporary mock for the approval prototype; replace with a real OTP challenge.
		if req.OTP != "123456" {
			return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "รหัส OTP ไม่ถูกต้อง")
		}
		if action == "reject" && comment == "" {
			comment = "ผอ.สำนักไม่อนุมัติคำร้อง"
		}
	}
	if (action == "return_revision" || action == "reject") && comment == "" {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุเหตุผล")
	}
	updateReq := requests.UpdateApplicationStatusRequest{Reason: comment}
	if isBureau {
		input := req.BureauDirectorReview
		if input == nil && !isStallReview {
			input = &models.BureauDirectorReview{Decision: action}
		}
		review, reviewErr := encodeBureauDirectorReview(input)
		if reviewErr != nil {
			return nil, reviewErr
		}
		if input.Decision != action {
			return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "ผลพิจารณาไม่ตรงกับคำสั่ง")
		}
		updateReq.BureauDirectorReviewJSON = review
	}
	if action == "approve" {
		signature, signatureType, signatureErr := decodeApprovalSignature(req.SignatureBase64)
		if signatureErr != nil {
			return nil, signatureErr
		}
		if isBureau {
			updateReq.BureauDirectorSignature, updateReq.BureauDirectorSignatureType = signature, signatureType
		} else {
			updateReq.MarketDirectorSignature, updateReq.MarketDirectorSignatureType = signature, signatureType
		}
	}
	switch action {
	case "return_revision":
		updateReq.Status = "revision_required"
	case "reject":
		updateReq.Status = "rejected"
	case "approve":
		if app.ApplicationState == "pending_market_director" {
			updateReq.Status = "pending_bureau_director"
		} else {
			updateReq.Status = "approved"
		}
	default:
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "การดำเนินการไม่ถูกต้อง")
	}
	if app.ApplicationState == "pending_market_director" {
		updateReq.MarketDirectorComment = comment
	} else {
		updateReq.BureauDirectorComment = comment
	}
	return u.UpdateStatus(id, userID, updateReq)
}

func decodeApprovalSignature(encoded string) ([]byte, string, error) {
	if len(encoded) == 0 || len(encoded) > 4*1024*1024 {
		return nil, "", responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาแนบลายเซ็นขนาดไม่เกิน 3 MB")
	}
	signature, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || len(signature) == 0 || len(signature) > 3*1024*1024 {
		return nil, "", responses.NewAppErr(http.StatusUnprocessableEntity, "ไฟล์ลายเซ็นไม่ถูกต้อง")
	}
	kind := http.DetectContentType(signature)
	if kind != "image/png" && kind != "image/jpeg" {
		return nil, "", responses.NewAppErr(http.StatusUnprocessableEntity, "รองรับลายเซ็น PNG หรือ JPEG เท่านั้น")
	}
	return signature, kind, nil
}
