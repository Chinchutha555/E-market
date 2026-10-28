package usecases

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"emarket-api/logs"
	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/contract/domains"
	"gorm.io/gorm"
)

type useCase struct {
	repo domains.Repository
}

func NewUseCase(repo domains.Repository) domains.UseCase {
	return &useCase{repo: repo}
}

func toContractResponse(c *models.Contract) responses.ContractResponse {
	return responses.ContractResponse{
		ID:              c.ID,
		ContractNumber:  c.ContractNumber,
		ApplicationID:   c.ApplicationID,
		StallID:         c.StallID,
		MarketID:        c.MarketID,
		TenantUserID:    c.TenantUserID,
		StallIdent:      c.StallIdent,
		TenantName:      c.TenantName,
		TenantPhone:     c.TenantPhone,
		TenantIDCard:    c.TenantIDCard,
		ShopName:        c.ShopName,
		StartDate:       c.StartDate,
		EndDate:         c.EndDate,
		MonthlyRental:   c.MonthlyRental,
		SecurityDeposit: c.SecurityDeposit,
		Status:          c.ContractState,
		CreatedAt:       c.CreatedAt,
	}
}

func (u *useCase) GetContracts(filter requests.ContractFilter) ([]responses.ContractResponse, error) {
	contracts, err := u.repo.FindAll(filter)
	if err != nil {
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	res := make([]responses.ContractResponse, len(contracts))
	for i, c := range contracts {
		res[i] = toContractResponse(&c)
	}
	return res, nil
}

func (u *useCase) GetContractByID(id int) (*responses.ContractResponse, error) {
	contract, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลสัญญา")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	res := toContractResponse(contract)
	return &res, nil
}

func (u *useCase) CreateContract(userID int, userRole string, req requests.CreateContractRequest) (*responses.ContractResponse, error) {
	stallIdent := strings.TrimSpace(req.StallIdent)
	if stallIdent == "" {
		stallIdent = strings.TrimSpace(req.StallNumber)
	}
	if stallIdent == "" && req.StallID > 0 {
		stallIdent = fmt.Sprintf("%d", req.StallID)
	}
	if stallIdent == "" {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุรหัสแผงค้า")
	}

	tenantName := strings.TrimSpace(req.TenantName)
	if tenantName == "" {
		return nil, responses.NewAppErr(http.StatusUnprocessableEntity, "กรุณาระบุชื่อผู้เช่า")
	}

	tenantPhone := strings.TrimSpace(req.TenantPhone)
	if tenantPhone == "" {
		tenantPhone = "081-234-5678"
	}

	shopName := strings.TrimSpace(req.ShopName)
	if shopName == "" {
		shopName = strings.TrimSpace(req.BusinessName)
	}
	if shopName == "" {
		shopName = tenantName
	}

	startDate := strings.TrimSpace(req.StartDate)
	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}

	endDate := strings.TrimSpace(req.EndDate)
	if endDate == "" {
		endDate = time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	}

	monthlyRental := req.MonthlyRental
	if monthlyRental <= 0 {
		monthlyRental = req.RentalFee
	}
	if monthlyRental <= 0 {
		monthlyRental = 2000
	}

	securityDeposit := req.SecurityDeposit
	if securityDeposit <= 0 {
		securityDeposit = monthlyRental * 2
	}

	contractNo := fmt.Sprintf("CTR-%s-%06d", time.Now().Format("20060102"), time.Now().UnixMicro()%1000000)

	var appID *int
	if req.ApplicationID > 0 {
		appID = &req.ApplicationID
	}

	var stallIDPtr *int
	if req.StallID > 0 {
		stallIDPtr = &req.StallID
	}

	contractState := "active"
	if userRole == "tenant" {
		contractState = "pending"
	} else if req.Status != "" {
		contractState = req.Status
	}

	contract := models.Contract{
		ContractNumber:  contractNo,
		ApplicationID:   appID,
		StallID:         stallIDPtr,
		StallIdent:      stallIdent,
		TenantName:      tenantName,
		TenantPhone:     tenantPhone,
		TenantIDCard:    req.TenantIDCard,
		ShopName:        shopName,
		StartDate:       startDate,
		EndDate:         endDate,
		MonthlyRental:   monthlyRental,
		SecurityDeposit: securityDeposit,
		ContractState:   contractState,
	}
	if userID > 0 {
		contract.CreatedBy = &userID
		contract.UpdatedBy = &userID
		if userRole == "tenant" {
			contract.TenantUserID = &userID
		}
	}

	if err := u.repo.Create(&contract); err != nil {
		if errors.Is(err, domains.ErrApplicationNotApproved) {
			return nil, responses.NewAppErr(http.StatusConflict, "คำร้องยังไม่ได้รับอนุมัติขั้นสุดท้าย หรือไม่ตรงกับแผงที่เลือก")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}

	res := toContractResponse(&contract)
	return &res, nil
}

func (u *useCase) ApproveContract(id int, userID int) (*responses.ContractResponse, error) {
	contract, err := u.repo.Approve(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลสัญญา")
		}
		logs.Error(err)
		return nil, responses.NewInternalServerError()
	}
	res := toContractResponse(contract)
	return &res, nil
}

func (u *useCase) RejectContract(id int, userID int, reason string) error {
	if err := u.repo.Reject(id, userID, reason); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลสัญญา")
		}
		logs.Error(err)
		return responses.NewInternalServerError()
	}
	return nil
}

func (u *useCase) TerminateContract(id int, userID int) error {
	if err := u.repo.Terminate(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.NewAppErr(http.StatusNotFound, "ไม่พบข้อมูลสัญญา")
		}
		logs.Error(err)
		return responses.NewInternalServerError()
	}
	return nil
}
