package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
)

// ChatRepository handles chat data operations
type ChatRepository struct {
	db *gorm.DB
}

// NewChatRepository creates a new chat repository
func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// FindAll returns all chat messages
func (r *ChatRepository) FindAll() ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Order("created_at desc").Find(&messages).Error
	return messages, err
}

// FindByID returns a chat message by ID
func (r *ChatRepository) FindByID(id uint) (*models.ChatMessage, error) {
	var message models.ChatMessage
	err := r.db.First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// FindByContactID returns all messages for a contact
func (r *ChatRepository) FindByContactID(contactID uint) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Where("contact_id = ?", contactID).Order("created_at desc").Find(&messages).Error
	return messages, err
}

// Create creates a new chat message
func (r *ChatRepository) Create(message *models.ChatMessage) error {
	return r.db.Create(message).Error
}

// Update updates a chat message
func (r *ChatRepository) Update(message *models.ChatMessage) error {
	return r.db.Save(message).Error
}

// Delete deletes a chat message
func (r *ChatRepository) Delete(id uint) error {
	return r.db.Delete(&models.ChatMessage{}, id).Error
}

// GetQuickReplies returns all quick replies
func (r *ChatRepository) GetQuickReplies() ([]models.QuickReply, error) {
	var replies []models.QuickReply
	err := r.db.Where("active = ?", true).Find(&replies).Error
	return replies, err
}

// GetChatLabels returns all chat labels
func (r *ChatRepository) GetChatLabels() ([]models.ChatLabel, error) {
	var labels []models.ChatLabel
	err := r.db.Find(&labels).Error
	return labels, err
}

// GetBroadcastTemplates returns all broadcast templates
func (r *ChatRepository) GetBroadcastTemplates() ([]models.BroadcastTemplate, error) {
	var templates []models.BroadcastTemplate
	err := r.db.Find(&templates).Error
	return templates, err
}
