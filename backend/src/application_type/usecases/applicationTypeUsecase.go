package usecases

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
	"emarket-api/src/application_type/domains"
)

type useCase struct {
	repo domains.Repository
}

func NewUseCase(repo domains.Repository) domains.UseCase {
	return &useCase{repo: repo}
}

func toResponse(m *models.ApplicationType) responses.ApplicationTypeResponse {
	return responses.ApplicationTypeResponse{
		ID:          m.ID,
		Code:        m.Code,
		Title:       m.Title,
		Description: m.Description,
		IsActive:    m.IsActive,
		SortOrder:   m.SortOrder,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (u *useCase) GetApplicationTypes(filter requests.ApplicationTypeFilter) ([]responses.ApplicationTypeResponse, error) {
	items, err := u.repo.FindAll(filter)
	if err != nil {
		return nil, responses.NewAppErr(http.StatusInternalServerError, "ไม่สามารถดึงข้อมูลประเภทคำร้องได้")
	}

	result := make([]responses.ApplicationTypeResponse, len(items))
	for i, item := range items {
		result[i] = toResponse(&item)
	}
	return result, nil
}

func (u *useCase) GetApplicationTypeByID(id uint) (*responses.ApplicationTypeResponse, error) {
	item, err := u.repo.FindByID(id)
	if err != nil {
		return nil, responses.NewAppErr(http.StatusInternalServerError, "ไม่สามารถค้นหาประเภทคำร้องได้")
	}
	if item == nil {
		return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบประเภทคำร้องนี้ในระบบ")
	}
	res := toResponse(item)
	return &res, nil
}

func (u *useCase) CreateApplicationType(req requests.CreateApplicationTypeRequest) (*responses.ApplicationTypeResponse, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, responses.NewAppErr(http.StatusBadRequest, "กรุณาระบุชื่อประเภทคำร้อง / เรื่องที่ขอ")
	}

	code := strings.TrimSpace(req.Code)
	if code == "" {
		// Generate slug or code
		code = fmt.Sprintf("app_type_%d", time.Now().UnixNano()%1000000)
	}

	// Check unique code
	existing, _ := u.repo.FindByCode(code)
	if existing != nil {
		code = fmt.Sprintf("%s_%d", code, time.Now().Unix()%10000)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	appType := models.ApplicationType{
		Title:       title,
		Code:        code,
		Description: strings.TrimSpace(req.Description),
		IsActive:    isActive,
		SortOrder:   req.SortOrder,
	}

	if err := u.repo.Create(&appType); err != nil {
		return nil, responses.NewAppErr(http.StatusInternalServerError, "เกิดข้อผิดพลาดในการบันทึกประเภทคำร้อง")
	}

	res := toResponse(&appType)
	return &res, nil
}

func (u *useCase) UpdateApplicationType(id uint, req requests.UpdateApplicationTypeRequest) (*responses.ApplicationTypeResponse, error) {
	item, err := u.repo.FindByID(id)
	if err != nil {
		return nil, responses.NewAppErr(http.StatusInternalServerError, "ไม่สามารถค้นหาประเภทคำร้องได้")
	}
	if item == nil {
		return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบประเภทคำร้องนี้ในระบบ")
	}

	if req.Title != "" {
		item.Title = strings.TrimSpace(req.Title)
	}
	if req.Code != "" {
		newCode := strings.TrimSpace(req.Code)
		if newCode != item.Code {
			existing, _ := u.repo.FindByCode(newCode)
			if existing != nil && existing.ID != item.ID {
				return nil, responses.NewAppErr(http.StatusBadRequest, "รหัสประเภทคำร้อง (Code) นี้ถูกใช้งานแล้ว")
			}
			item.Code = newCode
		}
	}
	item.Description = strings.TrimSpace(req.Description)
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		item.SortOrder = *req.SortOrder
	}

	if err := u.repo.Update(item); err != nil {
		return nil, responses.NewAppErr(http.StatusInternalServerError, "เกิดข้อผิดพลาดในการแก้ไขประเภทคำร้อง")
	}

	res := toResponse(item)
	return &res, nil
}

func (u *useCase) DeleteApplicationType(id uint) error {
	item, err := u.repo.FindByID(id)
	if err != nil {
		return responses.NewAppErr(http.StatusInternalServerError, "ไม่สามารถค้นหาประเภทคำร้องได้")
	}
	if item == nil {
		return responses.NewAppErr(http.StatusNotFound, "ไม่พบประเภทคำร้องนี้ในระบบ")
	}

	if err := u.repo.Delete(id); err != nil {
		return responses.NewAppErr(http.StatusInternalServerError, "เกิดข้อผิดพลาดในการลบประเภทคำร้อง")
	}
	return nil
}

func (u *useCase) ToggleActive(id uint) (*responses.ApplicationTypeResponse, error) {
	item, err := u.repo.ToggleActive(id)
	if err != nil {
		return nil, responses.NewAppErr(http.StatusInternalServerError, "เกิดข้อผิดพลาดในการเปลี่ยนสถานะ")
	}
	if item == nil {
		return nil, responses.NewAppErr(http.StatusNotFound, "ไม่พบประเภทคำร้องนี้ในระบบ")
	}
	res := toResponse(item)
	return &res, nil
}
