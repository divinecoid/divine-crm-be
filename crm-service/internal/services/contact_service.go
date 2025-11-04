package services

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"fmt"
	"time"
)

// ContactService handles contact business logic
type ContactService struct {
	repo   *repository.ContactRepository
	logger *utils.Logger
}

// NewContactService creates a new contact service
func NewContactService(repo *repository.ContactRepository, logger *utils.Logger) *ContactService {
	return &ContactService{
		repo:   repo,
		logger: logger,
	}
}

// GetAll returns all contacts
func (s *ContactService) GetAll() ([]models.Contact, error) {
	return s.repo.FindAll()
}

// GetByID returns a contact by ID
func (s *ContactService) GetByID(id uint) (*models.Contact, error) {
	return s.repo.FindByID(id)
}

// Create creates a new contact
func (s *ContactService) Create(contact *models.Contact) error {
	// Generate code if not provided
	if contact.Code == "" {
		count, _ := s.repo.Count()
		contact.Code = fmt.Sprintf("C%06d", count+1)
	}

	// Set timestamps
	contact.FirstContact = time.Now()
	contact.LastContact = time.Now()

	s.logger.Info("Creating new contact", "code", contact.Code, "name", contact.Name)
	return s.repo.Create(contact)
}

// Update updates a contact
func (s *ContactService) Update(contact *models.Contact) error {
	s.logger.Info("Updating contact", "id", contact.ID, "name", contact.Name)
	return s.repo.Update(contact)
}

// Delete deletes a contact
func (s *ContactService) Delete(id uint) error {
	s.logger.Info("Deleting contact", "id", id)
	return s.repo.Delete(id)
}

// GetOrCreateByChannelID gets or creates a contact by channel ID
func (s *ContactService) GetOrCreateByChannelID(channel, channelID, name string) (*models.Contact, error) {
	// Try to find existing contact
	contact, err := s.repo.FindByChannelID(channel, channelID)

	if contact == nil {
		// Contact not found, create new one
		s.logger.Info("Contact not found, creating new", "channel", channel, "channelID", channelID)

		count, _ := s.repo.Count()
		code := fmt.Sprintf("C%06d", count+1)

		contact = &models.Contact{
			Code:          code,
			Channel:       channel,
			ChannelID:     channelID,
			Name:          name,
			Temperature:   "Warm",
			FirstContact:  time.Now(),
			LastContact:   time.Now(),
			LastAgent:     "AI",
			LastAgentType: "Bot",
		}

		if err := s.repo.Create(contact); err != nil {
			s.logger.Error("Failed to create contact", "error", err)
			return nil, err
		}

		s.logger.Info("Created new contact", "code", contact.Code, "id", contact.ID)
		return contact, nil
	}

	// Contact found, check for errors
	if err != nil {
		s.logger.Error("Error finding contact", "error", err)
		return nil, err
	}

	// Update existing contact
	s.logger.Info("Found existing contact", "id", contact.ID, "code", contact.Code)
	contact.LastContact = time.Now()
	if name != "" && contact.Name != name {
		contact.Name = name
	}

	if err := s.repo.Update(contact); err != nil {
		s.logger.Warn("Failed to update contact timestamp", "error", err)
		// Don't fail, just warn
	}

	return contact, nil
}
