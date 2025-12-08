package services

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
)

type AIAgentService struct {
	repo   *repository.AIRepository
	logger *utils.Logger
}

func NewAIAgentService(repo *repository.AIRepository, logger *utils.Logger) *AIAgentService {
	return &AIAgentService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AIAgentService) GetAll() ([]models.AIAgent, error) {
	return s.repo.GetAllAgents()
}

func (s *AIAgentService) GetActive() (*models.AIAgent, error) {
	return s.repo.GetActiveAgent()
}

func (s *AIAgentService) GetByID(id uint) (*models.AIAgent, error) {
	return s.repo.GetAgentByID(id)
}

func (s *AIAgentService) Create(agent *models.AIAgent) error {
	s.logger.Info("Creating AI agent", "name", agent.Name)
	return s.repo.CreateAgent(agent)
}

func (s *AIAgentService) Update(agent *models.AIAgent) error {
	s.logger.Info("Updating AI agent", "id", agent.ID)
	return s.repo.UpdateAgent(agent)
}

func (s *AIAgentService) Delete(id uint) error {
	s.logger.Info("Deleting AI agent", "id", id)
	return s.repo.DeleteAgent(id)
}

func (s *AIAgentService) Activate(id uint) error {
	// Deactivate all agents first
	agents, err := s.repo.GetAllAgents()
	if err != nil {
		return err
	}

	for _, agent := range agents {
		agent.Active = false
		s.repo.UpdateAgent(&agent)
	}

	// Activate the selected agent
	agent, err := s.repo.GetAgentByID(id)
	if err != nil {
		return err
	}

	agent.Active = true
	s.logger.Info("Activating AI agent", "id", id, "name", agent.Name)
	return s.repo.UpdateAgent(agent)
}
