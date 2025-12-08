package services

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
)

type ChatLabelService struct {
	repo   *repository.ChatLabelRepository
	logger *utils.Logger
}

func NewChatLabelService(repo *repository.ChatLabelRepository, logger *utils.Logger) *ChatLabelService {
	return &ChatLabelService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ChatLabelService) GetAll() ([]models.ChatLabel, error) {
	return s.repo.FindAll()
}

func (s *ChatLabelService) GetByID(id uint) (*models.ChatLabel, error) {
	return s.repo.FindByID(id)
}

func (s *ChatLabelService) Create(label *models.ChatLabel) error {
	s.logger.Info("Creating chat label", "label", label.Label)
	return s.repo.Create(label)
}

func (s *ChatLabelService) Update(label *models.ChatLabel) error {
	s.logger.Info("Updating chat label", "id", label.ID)
	return s.repo.Update(label)
}

func (s *ChatLabelService) Delete(id uint) error {
	s.logger.Info("Deleting chat label", "id", id)
	return s.repo.Delete(id)
}
