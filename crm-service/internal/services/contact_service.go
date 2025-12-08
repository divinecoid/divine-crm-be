package services

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"fmt"
	"time"
)

type ContactService struct {
	repo   *repository.ContactRepository
	logger *utils.Logger
}

func NewContactService(repo *repository.ContactRepository, logger *utils.Logger) *ContactService {
	return &ContactService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ContactService) GetAll() ([]models.Contact, error) {
	return s.repo.FindAll()
}

func (s *ContactService) GetByID(id uint) (*models.Contact, error) {
	return s.repo.FindByID(id)
}

func (s *ContactService) GetByStatus(status string) ([]models.Contact, error) {
	return s.repo.FindByStatus(status)
}

func (s *ContactService) GetByTemperature(temp string) ([]models.Contact, error) {
	return s.repo.FindByTemperature(temp)
}

func (s *ContactService) Search(query string) ([]models.Contact, error) {
	return s.repo.Search(query)
}

func (s *ContactService) Create(contact *models.Contact) error {
	if contact.Code == "" {
		count, _ := s.repo.Count()
		contact.Code = fmt.Sprintf("C%06d", count+1)
	}

	contact.FirstContact = time.Now()
	contact.LastContact = time.Now()

	s.logger.Info("Creating new contact", "code", contact.Code, "name", contact.Name)
	return s.repo.Create(contact)
}

func (s *ContactService) Update(contact *models.Contact) error {
	s.logger.Info("Updating contact", "id", contact.ID, "name", contact.Name)
	return s.repo.Update(contact)
}

func (s *ContactService) Delete(id uint) error {
	s.logger.Info("Deleting contact", "id", id)
	return s.repo.Delete(id)
}

func (s *ContactService) GetOrCreateByChannelID(channel, channelID, name string) (*models.Contact, error) {
	contact, err := s.repo.FindByChannelID(channel, channelID)

	if contact == nil {
		s.logger.Info("Contact not found, creating new", "channel", channel, "channelID", channelID)

		count, _ := s.repo.Count()
		code := fmt.Sprintf("C%06d", count+1)

		contact = &models.Contact{
			Code:          code,
			Channel:       channel,
			ChannelID:     channelID,
			Name:          name,
			ContactStatus: "Leads",
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
	}

	return contact, nil
}

func (s *ContactService) UpdateTemperature(id uint, temperature string) error {
	contact, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	contact.Temperature = temperature
	return s.repo.Update(contact)
}

func (s *ContactService) UpdateStatus(id uint, status string) error {
	contact, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	contact.ContactStatus = status
	return s.repo.Update(contact)
}

// Statistics
func (s *ContactService) GetStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	total, _ := s.repo.Count()
	stats["total"] = total

	leads, _ := s.repo.CountByStatus("Leads")
	stats["leads"] = leads

	contacts, _ := s.repo.CountByStatus("Contact")
	stats["contacts"] = contacts

	return stats, nil
}
