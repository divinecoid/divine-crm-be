package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
)

// AIRepository handles AI data operations
type AIRepository struct {
	db *gorm.DB
}

// NewAIRepository creates a new AI repository
func NewAIRepository(db *gorm.DB) *AIRepository {
	return &AIRepository{db: db}
}

// GetActiveAgent returns the first active AI agent
func (r *AIRepository) GetActiveAgent() (*models.AIAgent, error) {
	var agent models.AIAgent
	err := r.db.Where("active = ?", true).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetAgentByID returns an AI agent by ID
func (r *AIRepository) GetAgentByID(id uint) (*models.AIAgent, error) {
	var agent models.AIAgent
	err := r.db.First(&agent, id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetConfiguration returns AI configuration by engine name
func (r *AIRepository) GetConfiguration(engine string) (*models.AIConfiguration, error) {
	var config models.AIConfiguration
	err := r.db.Where("ai_engine = ? AND active = ?", engine, true).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAllConfigurations returns all AI configurations
func (r *AIRepository) GetAllConfigurations() ([]models.AIConfiguration, error) {
	var configs []models.AIConfiguration
	err := r.db.Find(&configs).Error
	return configs, err
}

// GetAllAgents returns all AI agents
func (r *AIRepository) GetAllAgents() ([]models.AIAgent, error) {
	var agents []models.AIAgent
	err := r.db.Find(&agents).Error
	return agents, err
}

// CreateConfiguration creates a new AI configuration
func (r *AIRepository) CreateConfiguration(config *models.AIConfiguration) error {
	return r.db.Create(config).Error
}

// UpdateConfiguration updates an AI configuration
func (r *AIRepository) UpdateConfiguration(config *models.AIConfiguration) error {
	return r.db.Save(config).Error
}

// CreateAgent creates a new AI agent
func (r *AIRepository) CreateAgent(agent *models.AIAgent) error {
	return r.db.Create(agent).Error
}

// UpdateAgent updates an AI agent
func (r *AIRepository) UpdateAgent(agent *models.AIAgent) error {
	return r.db.Save(agent).Error
}

// GetTokenBalances returns all token balances
func (r *AIRepository) GetTokenBalances() ([]models.TokenBalance, error) {
	var balances []models.TokenBalance
	err := r.db.Find(&balances).Error
	return balances, err
}
