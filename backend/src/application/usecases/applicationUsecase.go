package usecases

import (
	"encoding/base64"
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
	if periodMonths <= 0 {
		periodMonths = 12
	}
	startDate := a.RequestedStartDate
	if startDate == "" {
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

	return responses.ApplicationResponse{
		ID:                    a.ID,
		ApplicationNo:         a.ApplicationNo,
		ApplicationNumber:     a.ApplicationNo,
		ApplicationType:       appType,
		Subject:               subject,
		BookingID:             a.BookingID,
		StallID:               a.StallID,
		MarketID:              a.MarketID,
		MarketName:            a.MarketName,
		StallIdent:            a.StallIdent,
		StallNumber:           a.StallIdent,
		ApplicantName:         a.ApplicantName,
		TenantName:            a.ApplicantName,
		ApplicantPhone:        a.ApplicantPhone,
		TenantPhone:           a.ApplicantPhone,
		ApplicantIDCard:       a.ApplicantIDCard,
		ApplicantAge:          a.ApplicantAge,
		HasIDCardCopy:         len(a.IDCardCopy) > 0,
		ApplicantEmail:        a.ApplicantEmail,
		ApplicantAddress:      a.ApplicantAddress,
		Details:               a.Details,
		ShopName:              a.ShopName,
		ProductCategory:       a.ProductCategory,
		ProductType:           a.ProductCategory,
		ProposedRentType:      a.ProposedRentType,
		RequestedStartDate:    startDate,
		LeaseTermMonths:       periodMonths,
		RequestedPeriodMonths: periodMonths,
		MonthlyRate:           monthlyRate,
		StallMonthlyRate:      monthlyRate,
		Status:                a.ApplicationState,
		OfficerComment:        a.OfficerComment,
		MarketDirectorComment: a.MarketDirectorComment,
		BureauDirectorComment: a.BureauDirectorComment,
		ReviewedAt:            a.ReviewedAt,
		ReviewNote:            a.ReviewNote,
		CreatedAt:             a.CreatedAt,
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
	if startDate == "" {
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
			_ = u.repo.UpdateStallStatus(app.StallIdent, "occupied")
		} else if req.Status == "rejected" {
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

	if app.ApplicationState == "pending_officer" {
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
		return u.UpdateStatus(id, userID, requests.UpdateApplicationStatusRequest{Status: status, Reason: comment, OfficerComment: comment})
	}
	if app.ApplicationState != "pending_market_director" && app.ApplicationState != "pending_bureau_director" {
		return nil, responses.NewAppErr(http.StatusConflict, "คำร้องไม่ได้อยู่ในขั้นที่สามารถพิจารณาได้")
	}
	action := strings.ToLower(strings.TrimSpace(req.Action))
	if (action == "return_revision" || action == "reject") && comment == "" {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุหมายเหตุ")
	}
	updateReq := requests.UpdateApplicationStatusRequest{Reason: comment}
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
