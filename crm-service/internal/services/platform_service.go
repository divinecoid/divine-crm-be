package services

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
)

type PlatformService struct {
	repo   *repository.PlatformRepository
	logger *utils.Logger
}

func NewPlatformService(repo *repository.PlatformRepository, logger *utils.Logger) *PlatformService {
	return &PlatformService{
		repo:   repo,
		logger: logger,
	}
}

func (s *PlatformService) GetAll() ([]models.ConnectedPlatform, error) {
	return s.repo.FindAll()
}

func (s *PlatformService) GetActive() ([]models.ConnectedPlatform, error) {
	return s.repo.FindActive()
}

func (s *PlatformService) GetByID(id uint) (*models.ConnectedPlatform, error) {
	return s.repo.FindByID(id)
}

func (s *PlatformService) Create(platform *models.ConnectedPlatform) error {
	s.logger.Info("Creating connected platform", "platform", platform.Platform)
	return s.repo.Create(platform)
}

func (s *PlatformService) Update(platform *models.ConnectedPlatform) error {
	s.logger.Info("Updating connected platform", "id", platform.ID)
	return s.repo.Update(platform)
}

func (s *PlatformService) Delete(id uint) error {
	s.logger.Info("Deleting connected platform", "id", id)
	return s.repo.Delete(id)
}
