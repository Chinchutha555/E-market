package domains

import (
	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/responses"
)

type Repository interface {
	FindAll(filter requests.ApplicationTypeFilter) ([]models.ApplicationType, error)
	FindByID(id uint) (*models.ApplicationType, error)
	FindByCode(code string) (*models.ApplicationType, error)
	Create(appType *models.ApplicationType) error
	Update(appType *models.ApplicationType) error
	Delete(id uint) error
	ToggleActive(id uint) (*models.ApplicationType, error)
}

type UseCase interface {
	GetApplicationTypes(filter requests.ApplicationTypeFilter) ([]responses.ApplicationTypeResponse, error)
	GetApplicationTypeByID(id uint) (*responses.ApplicationTypeResponse, error)
	CreateApplicationType(req requests.CreateApplicationTypeRequest) (*responses.ApplicationTypeResponse, error)
	UpdateApplicationType(id uint, req requests.UpdateApplicationTypeRequest) (*responses.ApplicationTypeResponse, error)
	DeleteApplicationType(id uint) error
	ToggleActive(id uint) (*responses.ApplicationTypeResponse, error)
}
