package services

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"strings"
	"time"
)

type BroadcastService struct {
	repo             *repository.BroadcastRepository
	contactRepo      *repository.ContactRepository
	whatsappService  *WhatsAppService
	instagramService *InstagramService
	telegramService  *TelegramService
	logger           *utils.Logger
}

func NewBroadcastService(
	repo *repository.BroadcastRepository,
	contactRepo *repository.ContactRepository,
	whatsappService *WhatsAppService,
	instagramService *InstagramService,
	telegramService *TelegramService,
	logger *utils.Logger,
) *BroadcastService {
	return &BroadcastService{
		repo:             repo,
		contactRepo:      contactRepo,
		whatsappService:  whatsappService,
		instagramService: instagramService,
		telegramService:  telegramService,
		logger:           logger,
	}
}

// Template Management
func (s *BroadcastService) GetAllTemplates() ([]models.BroadcastTemplate, error) {
	return s.repo.FindAllTemplates()
}

func (s *BroadcastService) GetTemplateByID(id uint) (*models.BroadcastTemplate, error) {
	return s.repo.FindTemplateByID(id)
}

func (s *BroadcastService) CreateTemplate(template *models.BroadcastTemplate) error {
	s.logger.Info("Creating broadcast template", "name", template.Name)
	return s.repo.CreateTemplate(template)
}

func (s *BroadcastService) UpdateTemplate(template *models.BroadcastTemplate) error {
	s.logger.Info("Updating broadcast template", "id", template.ID)
	return s.repo.UpdateTemplate(template)
}

func (s *BroadcastService) DeleteTemplate(id uint) error {
	s.logger.Info("Deleting broadcast template", "id", id)
	return s.repo.DeleteTemplate(id)
}

// Broadcast Execution
func (s *BroadcastService) SendBroadcast(templateID uint, sentBy string) error {
	template, err := s.repo.FindTemplateByID(templateID)
	if err != nil {
		return err
	}

	// Get recipients based on channel
	var contacts []models.Contact
	if template.Channel == "All" {
		contacts, err = s.contactRepo.FindAll()
	} else {
		contacts, err = s.contactRepo.Search(template.Channel)
	}

	if err != nil {
		return err
	}

	// Create broadcast history
	history := &models.BroadcastHistory{
		TemplateID: templateID,
		SentTo:     len(contacts),
		Successful: 0,
		Failed:     0,
		Status:     "Processing",
		SentBy:     sentBy,
		CreatedAt:  time.Now(),
	}

	s.repo.CreateHistory(history)

	// Send to each contact
	go s.processBroadcast(history, template, contacts)

	return nil
}

func (s *BroadcastService) processBroadcast(history *models.BroadcastHistory, template *models.BroadcastTemplate, contacts []models.Contact) {
	successful := 0
	failed := 0

	for _, contact := range contacts {
		// Personalize message
		message := s.personalizeMessage(template.Content, contact)

		// Send based on channel
		var err error
		switch contact.Channel {
		case "WhatsApp":
			err = s.whatsappService.SendMessage(contact.ChannelID, message)
		case "Instagram":
			err = s.instagramService.SendMessage(contact.ChannelID, message)
		case "Telegram":
			err = s.telegramService.SendMessage(contact.ChannelID, message)
		}

		if err != nil {
			failed++
			s.logger.Error("Failed to send broadcast", "contact", contact.Name, "error", err)
		} else {
			successful++
		}

		// Small delay to avoid rate limits
		time.Sleep(100 * time.Millisecond)
	}

	// Update history
	now := time.Now()
	history.Successful = successful
	history.Failed = failed
	history.Status = "Completed"
	history.CompletedAt = &now

	s.repo.UpdateHistory(history)
	s.logger.Info("Broadcast completed", "successful", successful, "failed", failed)
}

func (s *BroadcastService) personalizeMessage(template string, contact models.Contact) string {
	message := template
	message = strings.ReplaceAll(message, "{name}", contact.Name)
	message = strings.ReplaceAll(message, "{code}", contact.Code)
	return message
}

// History
func (s *BroadcastService) GetHistory() ([]models.BroadcastHistory, error) {
	return s.repo.FindAllHistory()
}
