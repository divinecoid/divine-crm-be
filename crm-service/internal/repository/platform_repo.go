package repository

import (
	"divine-crm/internal/models"
	"errors"
	"gorm.io/gorm"
)

type PlatformRepository struct {
	db *gorm.DB
}

func NewPlatformRepository(db *gorm.DB) *PlatformRepository {
	return &PlatformRepository{db: db}
}

func (r *PlatformRepository) FindAll() ([]models.ConnectedPlatform, error) {
	var platforms []models.ConnectedPlatform
	err := r.db.Order("platform ASC").Find(&platforms).Error
	return platforms, err
}

func (r *PlatformRepository) FindActive() ([]models.ConnectedPlatform, error) {
	var platforms []models.ConnectedPlatform
	err := r.db.Where("active = ?", true).Order("platform ASC").Find(&platforms).Error
	return platforms, err
}

func (r *PlatformRepository) FindByID(id uint) (*models.ConnectedPlatform, error) {
	var platform models.ConnectedPlatform
	err := r.db.First(&platform, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &platform, err
}

func (r *PlatformRepository) GetByPlatform(platform string) (*models.ConnectedPlatform, error) {
	var p models.ConnectedPlatform
	err := r.db.Where("platform = ?", platform).First(&p).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (r *PlatformRepository) Create(platform *models.ConnectedPlatform) error {
	return r.db.Create(platform).Error
}

func (r *PlatformRepository) Update(platform *models.ConnectedPlatform) error {
	return r.db.Save(platform).Error
}

func (r *PlatformRepository) Delete(id uint) error {
	return r.db.Delete(&models.ConnectedPlatform{}, id).Error
}
