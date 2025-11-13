package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
)

type BroadcastRepository struct {
	db *gorm.DB
}

func NewBroadcastRepository(db *gorm.DB) *BroadcastRepository {
	return &BroadcastRepository{db: db}
}

// Templates
func (r *BroadcastRepository) FindAllTemplates() ([]models.BroadcastTemplate, error) {
	var templates []models.BroadcastTemplate
	err := r.db.Order("created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *BroadcastRepository) FindTemplateByID(id uint) (*models.BroadcastTemplate, error) {
	var template models.BroadcastTemplate
	err := r.db.First(&template, id).Error
	return &template, err
}

func (r *BroadcastRepository) CreateTemplate(template *models.BroadcastTemplate) error {
	return r.db.Create(template).Error
}

func (r *BroadcastRepository) UpdateTemplate(template *models.BroadcastTemplate) error {
	return r.db.Save(template).Error
}

func (r *BroadcastRepository) DeleteTemplate(id uint) error {
	return r.db.Delete(&models.BroadcastTemplate{}, id).Error
}

// History
func (r *BroadcastRepository) FindAllHistory() ([]models.BroadcastHistory, error) {
	var history []models.BroadcastHistory
	err := r.db.Preload("Template").Order("created_at DESC").Find(&history).Error
	return history, err
}

func (r *BroadcastRepository) CreateHistory(history *models.BroadcastHistory) error {
	return r.db.Create(history).Error
}

func (r *BroadcastRepository) UpdateHistory(history *models.BroadcastHistory) error {
	return r.db.Save(history).Error
}
