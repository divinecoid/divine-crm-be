package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) FindAll() ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Preload("Contact").Order("created_at DESC").Find(&messages).Error
	return messages, err
}

func (r *ChatRepository) FindByID(id uint) (*models.ChatMessage, error) {
	var message models.ChatMessage
	err := r.db.Preload("Contact").First(&message, id).Error
	return &message, err
}

func (r *ChatRepository) FindByStatus(status string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Preload("Contact").
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

func (r *ChatRepository) FindByContactID(contactID uint) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Preload("Contact").
		Where("contact_id = ?", contactID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *ChatRepository) FindByChannel(channel string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Preload("Contact").
		Where("channel = ?", channel).
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

func (r *ChatRepository) Create(message *models.ChatMessage) error {
	return r.db.Create(message).Error
}

func (r *ChatRepository) Update(message *models.ChatMessage) error {
	return r.db.Save(message).Error
}

func (r *ChatRepository) Delete(id uint) error {
	return r.db.Delete(&models.ChatMessage{}, id).Error
}

// Statistics methods
func (r *ChatRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.ChatMessage{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *ChatRepository) CountByChannel(channel string) (int64, error) {
	var count int64
	err := r.db.Model(&models.ChatMessage{}).Where("channel = ?", channel).Count(&count).Error
	return count, err
}

func (r *ChatRepository) GetTotalTokens() (int64, error) {
	var total int64
	err := r.db.Model(&models.ChatMessage{}).Select("COALESCE(SUM(tokens_used), 0)").Scan(&total).Error
	return total, err
}

func (r *ChatRepository) GetTodayCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.ChatMessage{}).
		Where("DATE(created_at) = CURRENT_DATE").
		Count(&count).Error
	return count, err
}
