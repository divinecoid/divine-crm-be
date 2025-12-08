package services

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"strings"
)

type QuickReplyService struct {
	repo   *repository.QuickReplyRepository
	logger *utils.Logger
}

func NewQuickReplyService(repo *repository.QuickReplyRepository, logger *utils.Logger) *QuickReplyService {
	return &QuickReplyService{
		repo:   repo,
		logger: logger,
	}
}

func (s *QuickReplyService) GetAll() ([]models.QuickReply, error) {
	return s.repo.FindAll()
}

func (s *QuickReplyService) Create(reply *models.QuickReply) error {
	s.logger.Info("Creating quick reply", "trigger", reply.Trigger)
	return s.repo.Create(reply)
}

func (s *QuickReplyService) Update(reply *models.QuickReply) error {
	s.logger.Info("Updating quick reply", "id", reply.ID)
	return s.repo.Update(reply)
}

func (s *QuickReplyService) Delete(id uint) error {
	s.logger.Info("Deleting quick reply", "id", id)
	return s.repo.Delete(id)
}

// Check if message matches any quick reply
func (s *QuickReplyService) CheckQuickReply(message string) (*models.QuickReply, error) {
	messageLower := strings.ToLower(strings.TrimSpace(message))

	replies, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, reply := range replies {
		triggerLower := strings.ToLower(reply.Trigger)
		if strings.Contains(messageLower, triggerLower) {
			return &reply, nil
		}
	}

	return nil, nil
}
