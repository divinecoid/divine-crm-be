package services

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
)

type AIConfigService struct {
	repo   *repository.AIRepository
	logger *utils.Logger
}

func NewAIConfigService(repo *repository.AIRepository, logger *utils.Logger) *AIConfigService {
	return &AIConfigService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AIConfigService) GetAll() ([]models.AIConfiguration, error) {
	return s.repo.GetAllConfigurations()
}

func (s *AIConfigService) GetActive() (*models.AIConfiguration, error) {
	configs, err := s.repo.GetAllConfigurations()
	if err != nil {
		return nil, err
	}

	for _, config := range configs {
		if config.Active {
			return &config, nil
		}
	}

	return nil, nil
}

func (s *AIConfigService) GetByID(id uint) (*models.AIConfiguration, error) {
	// Implement based on repository method
	configs, err := s.repo.GetAllConfigurations()
	if err != nil {
		return nil, err
	}

	for _, config := range configs {
		if config.ID == id {
			return &config, nil
		}
	}

	return nil, nil
}

func (s *AIConfigService) Create(config *models.AIConfiguration) error {
	s.logger.Info("Creating AI configuration", "engine", config.AIEngine)
	return s.repo.CreateConfiguration(config)
}

func (s *AIConfigService) Update(config *models.AIConfiguration) error {
	s.logger.Info("Updating AI configuration", "id", config.ID)
	return s.repo.UpdateConfiguration(config)
}

func (s *AIConfigService) Delete(id uint) error {
	s.logger.Info("Deleting AI configuration", "id", id)
	return s.repo.DeleteConfiguration(id)
}
