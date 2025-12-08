package services

import (
	"divine-crm/internal/config"
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type HumanAgentService struct {
	repo   *repository.HumanAgentRepository
	logger *utils.Logger
	config *config.Config
}

func NewHumanAgentService(repo *repository.HumanAgentRepository, logger *utils.Logger, cfg *config.Config) *HumanAgentService {
	return &HumanAgentService{
		repo:   repo,
		logger: logger,
		config: cfg,
	}
}

func (s *HumanAgentService) GetAll() ([]models.HumanAgent, error) {
	return s.repo.FindAll()
}

func (s *HumanAgentService) GetByID(id uint) (*models.HumanAgent, error) {
	return s.repo.FindByID(id)
}

func (s *HumanAgentService) GetActive() ([]models.HumanAgent, error) {
	return s.repo.FindActive()
}

func (s *HumanAgentService) Create(agent *models.HumanAgent) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(agent.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	agent.Password = string(hashedPassword)

	s.logger.Info("Creating human agent", "username", agent.Username)
	return s.repo.Create(agent)
}

func (s *HumanAgentService) Update(agent *models.HumanAgent) error {
	s.logger.Info("Updating human agent", "id", agent.ID)
	return s.repo.Update(agent)
}

func (s *HumanAgentService) Delete(id uint) error {
	s.logger.Info("Deleting human agent", "id", id)
	return s.repo.Delete(id)
}

func (s *HumanAgentService) RevokeAccess(id uint) error {
	s.logger.Info("Revoking access for agent", "id", id)
	return s.repo.RevokeAccess(id)
}

func (s *HumanAgentService) ResetPassword(id uint, newPassword string) error {
	agent, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	agent.Password = string(hashedPassword)
	s.logger.Info("Resetting password for agent", "id", id)
	return s.repo.Update(agent)
}

func (s *HumanAgentService) Authenticate(username, password string) (*models.HumanAgent, error) {
	agent, err := s.repo.FindByUsername(username)
	if err != nil || agent == nil {
		return nil, errors.New("invalid username or password")
	}

	if !agent.Active {
		return nil, errors.New("account is inactive")
	}

	err = bcrypt.CompareHashAndPassword([]byte(agent.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Update login time
	s.repo.UpdateLoginTime(agent.ID)

	return agent, nil
}

func (s *HumanAgentService) GenerateToken(agent *models.HumanAgent) (string, error) {
	return utils.GenerateToken(
		agent.ID,
		agent.Email,
		agent.Role,
		s.config.JWT.Secret,
		s.config.JWT.ExpirationHours,
	)
}
