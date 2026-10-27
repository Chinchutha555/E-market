package repositories

import (
	"errors"
	"strings"

	"emarket-api/models"
	"emarket-api/requests"
	"emarket-api/src/application_type/domains"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domains.Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(filter requests.ApplicationTypeFilter) ([]models.ApplicationType, error) {
	var items []models.ApplicationType
	query := r.db.Model(&models.ApplicationType{})

	if filter.ActiveOnly {
		query = query.Where("is_active = ?", true)
	}

	if filter.Search != "" {
		s := "%" + strings.TrimSpace(filter.Search) + "%"
		query = query.Where("title LIKE ? OR code LIKE ? OR description LIKE ?", s, s, s)
	}

	err := query.Order("sort_order ASC, id ASC").Find(&items).Error
	return items, err
}

func (r *repository) FindByID(id uint) (*models.ApplicationType, error) {
	var item models.ApplicationType
	if err := r.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) FindByCode(code string) (*models.ApplicationType, error) {
	var item models.ApplicationType
	if err := r.db.Where("code = ?", code).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *repository) Create(appType *models.ApplicationType) error {
	return r.db.Create(appType).Error
}

func (r *repository) Update(appType *models.ApplicationType) error {
	return r.db.Save(appType).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.ApplicationType{}, id).Error
}

func (r *repository) ToggleActive(id uint) (*models.ApplicationType, error) {
	item, err := r.FindByID(id)
	if err != nil || item == nil {
		return nil, err
	}
	item.IsActive = !item.IsActive
	if err := r.db.Save(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}
