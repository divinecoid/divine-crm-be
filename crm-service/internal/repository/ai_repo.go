package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
)

type AIRepository struct {
	db *gorm.DB
}

func (r *AIRepository) GetTokenBalances() (any, any) {
	panic("unimplemented")
}

func NewAIRepository(db *gorm.DB) *AIRepository {
	return &AIRepository{db: db}
}

// ==================== AI CONFIGURATIONS ====================

func (r *AIRepository) GetAllConfigurations() ([]models.AIConfiguration, error) {
	var configs []models.AIConfiguration
	err := r.db.Order("ai_engine ASC").Find(&configs).Error
	return configs, err
}

func (r *AIRepository) GetConfiguration(engine string) (*models.AIConfiguration, error) {
	var config models.AIConfiguration
	err := r.db.Where("ai_engine = ?", engine).First(&config).Error
	return &config, err
}

func (r *AIRepository) CreateConfiguration(config *models.AIConfiguration) error {
	return r.db.Create(config).Error
}

func (r *AIRepository) UpdateConfiguration(config *models.AIConfiguration) error {
	return r.db.Save(config).Error
}

func (r *AIRepository) DeleteConfiguration(id uint) error {
	return r.db.Delete(&models.AIConfiguration{}, id).Error
}

// ==================== AI AGENTS ====================

func (r *AIRepository) GetAllAgents() ([]models.AIAgent, error) {
	var agents []models.AIAgent
	err := r.db.Order("name ASC").Find(&agents).Error
	return agents, err
}

func (r *AIRepository) GetActiveAgent() (*models.AIAgent, error) {
	var agent models.AIAgent
	err := r.db.Where("active = ?", true).First(&agent).Error
	return &agent, err
}

func (r *AIRepository) GetAgentByID(id uint) (*models.AIAgent, error) {
	var agent models.AIAgent
	err := r.db.First(&agent, id).Error
	return &agent, err
}

func (r *AIRepository) CreateAgent(agent *models.AIAgent) error {
	return r.db.Create(agent).Error
}

func (r *AIRepository) UpdateAgent(agent *models.AIAgent) error {
	return r.db.Save(agent).Error
}

func (r *AIRepository) DeleteAgent(id uint) error {
	return r.db.Delete(&models.AIAgent{}, id).Error
}
