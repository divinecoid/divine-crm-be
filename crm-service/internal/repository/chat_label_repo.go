package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
)

type ChatLabelRepository struct {
	db *gorm.DB
}

func NewChatLabelRepository(db *gorm.DB) *ChatLabelRepository {
	return &ChatLabelRepository{db: db}
}

func (r *ChatLabelRepository) FindAll() ([]models.ChatLabel, error) {
	var labels []models.ChatLabel
	err := r.db.Order("label ASC").Find(&labels).Error
	return labels, err
}

func (r *ChatLabelRepository) FindByID(id uint) (*models.ChatLabel, error) {
	var label models.ChatLabel
	err := r.db.First(&label, id).Error
	return &label, err
}

func (r *ChatLabelRepository) FindByLabel(label string) (*models.ChatLabel, error) {
	var chatLabel models.ChatLabel
	err := r.db.Where("label = ?", label).First(&chatLabel).Error
	return &chatLabel, err
}

func (r *ChatLabelRepository) Create(label *models.ChatLabel) error {
	return r.db.Create(label).Error
}

func (r *ChatLabelRepository) Update(label *models.ChatLabel) error {
	return r.db.Save(label).Error
}

func (r *ChatLabelRepository) Delete(id uint) error {
	return r.db.Delete(&models.ChatLabel{}, id).Error
}
