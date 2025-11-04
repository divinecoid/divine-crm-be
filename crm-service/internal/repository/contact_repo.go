package repository

import (
	"divine-crm/internal/models"
	"errors"
	"gorm.io/gorm"
)

// ContactRepository handles contact data operations
type ContactRepository struct {
	db *gorm.DB
}

// NewContactRepository creates a new contact repository
func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

// FindAll returns all contacts
func (r *ContactRepository) FindAll() ([]models.Contact, error) {
	var contacts []models.Contact
	err := r.db.Order("created_at desc").Find(&contacts).Error
	return contacts, err
}

// FindByID returns a contact by ID
func (r *ContactRepository) FindByID(id uint) (*models.Contact, error) {
	var contact models.Contact
	err := r.db.First(&contact, id).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// FindByChannelID returns a contact by channel and channel ID
func (r *ContactRepository) FindByChannelID(channel, channelID string) (*models.Contact, error) {
	var contact models.Contact
	err := r.db.Where("channel = ? AND channel_id = ?", channel, channelID).First(&contact).Error
	if err != nil {
		// ✅ Check if it's "not found" error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil (not an error!)
		}
		return nil, err // Return actual error
	}

	return &contact, nil
}

// Create creates a new contact
func (r *ContactRepository) Create(contact *models.Contact) error {
	return r.db.Create(contact).Error
}

// Update updates a contact
func (r *ContactRepository) Update(contact *models.Contact) error {
	return r.db.Save(contact).Error
}

// Delete deletes a contact
func (r *ContactRepository) Delete(id uint) error {
	return r.db.Delete(&models.Contact{}, id).Error
}

// Count returns total contacts
func (r *ContactRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Contact{}).Count(&count).Error
	return count, err
}
