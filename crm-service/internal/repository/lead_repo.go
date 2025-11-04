package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
)

// LeadRepository handles lead data operations
type LeadRepository struct {
	db *gorm.DB
}

// NewLeadRepository creates a new lead repository
func NewLeadRepository(db *gorm.DB) *LeadRepository {
	return &LeadRepository{db: db}
}

// FindAll returns all leads
func (r *LeadRepository) FindAll() ([]models.Lead, error) {
	var leads []models.Lead
	err := r.db.Order("created_at desc").Find(&leads).Error
	return leads, err
}

// FindByID returns a lead by ID
func (r *LeadRepository) FindByID(id uint) (*models.Lead, error) {
	var lead models.Lead
	err := r.db.First(&lead, id).Error
	if err != nil {
		return nil, err
	}
	return &lead, nil
}

// Create creates a new lead
func (r *LeadRepository) Create(lead *models.Lead) error {
	return r.db.Create(lead).Error
}

// Update updates a lead
func (r *LeadRepository) Update(lead *models.Lead) error {
	return r.db.Save(lead).Error
}

// Delete deletes a lead
func (r *LeadRepository) Delete(id uint) error {
	return r.db.Delete(&models.Lead{}, id).Error
}
