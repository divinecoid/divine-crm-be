package repository

import (
	"divine-crm/internal/models"
	"errors"
	"gorm.io/gorm"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) FindAll() ([]models.Contact, error) {
	var contacts []models.Contact
	err := r.db.Order("last_contact DESC").Find(&contacts).Error
	return contacts, err
}

func (r *ContactRepository) FindByID(id uint) (*models.Contact, error) {
	var contact models.Contact
	err := r.db.First(&contact, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &contact, err
}

func (r *ContactRepository) FindByChannelID(channel, channelID string) (*models.Contact, error) {
	var contact models.Contact
	err := r.db.Where("channel = ? AND channel_id = ?", channel, channelID).First(&contact).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &contact, err
}

func (r *ContactRepository) FindByStatus(status string) ([]models.Contact, error) {
	var contacts []models.Contact
	err := r.db.Where("contact_status = ?", status).Order("last_contact DESC").Find(&contacts).Error
	return contacts, err
}

func (r *ContactRepository) FindByTemperature(temp string) ([]models.Contact, error) {
	var contacts []models.Contact
	err := r.db.Where("temperature = ?", temp).Order("last_contact DESC").Find(&contacts).Error
	return contacts, err
}

func (r *ContactRepository) Search(query string) ([]models.Contact, error) {
	var contacts []models.Contact
	searchPattern := "%" + query + "%"
	err := r.db.Where("name ILIKE ? OR code ILIKE ? OR channel_id ILIKE ?",
		searchPattern, searchPattern, searchPattern).
		Order("last_contact DESC").
		Find(&contacts).Error
	return contacts, err
}

func (r *ContactRepository) Create(contact *models.Contact) error {
	return r.db.Create(contact).Error
}

func (r *ContactRepository) Update(contact *models.Contact) error {
	return r.db.Save(contact).Error
}

func (r *ContactRepository) Delete(id uint) error {
	return r.db.Delete(&models.Contact{}, id).Error
}

func (r *ContactRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Contact{}).Count(&count).Error
	return count, err
}

func (r *ContactRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Contact{}).Where("contact_status = ?", status).Count(&count).Error
	return count, err
}
