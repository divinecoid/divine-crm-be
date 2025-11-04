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

// GetByPlatform gets platform configuration by platform name
func (r *PlatformRepository) GetByPlatform(platform string) (*models.ConnectedPlatform, error) {
	var p models.ConnectedPlatform
	err := r.db.Where("platform = ?", platform).First(&p).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is ok
		}
		return nil, err
	}

	return &p, nil
}

// Create creates a new platform
func (r *PlatformRepository) Create(platform *models.ConnectedPlatform) error {
	return r.db.Create(platform).Error
}

// Update updates a platform
func (r *PlatformRepository) Update(platform *models.ConnectedPlatform) error {
	return r.db.Save(platform).Error
}
