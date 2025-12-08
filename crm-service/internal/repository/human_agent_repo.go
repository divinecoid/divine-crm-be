package repository

import (
	"divine-crm/internal/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

type HumanAgentRepository struct {
	db *gorm.DB
}

func NewHumanAgentRepository(db *gorm.DB) *HumanAgentRepository {
	return &HumanAgentRepository{db: db}
}

func (r *HumanAgentRepository) FindAll() ([]models.HumanAgent, error) {
	var agents []models.HumanAgent
	err := r.db.Order("username ASC").Find(&agents).Error
	return agents, err
}

func (r *HumanAgentRepository) FindByID(id uint) (*models.HumanAgent, error) {
	var agent models.HumanAgent
	err := r.db.First(&agent, id).Error
	return &agent, err
}

func (r *HumanAgentRepository) FindByUsername(username string) (*models.HumanAgent, error) {
	var agent models.HumanAgent
	err := r.db.Where("username = ?", username).First(&agent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &agent, err
}

func (r *HumanAgentRepository) FindActive() ([]models.HumanAgent, error) {
	var agents []models.HumanAgent
	err := r.db.Where("active = ?", true).Order("username ASC").Find(&agents).Error
	return agents, err
}

func (r *HumanAgentRepository) Create(agent *models.HumanAgent) error {
	return r.db.Create(agent).Error
}

func (r *HumanAgentRepository) Update(agent *models.HumanAgent) error {
	return r.db.Save(agent).Error
}

func (r *HumanAgentRepository) UpdateLoginTime(id uint) error {
	now := time.Now()
	return r.db.Model(&models.HumanAgent{}).Where("id = ?", id).Update("latest_login", now).Error
}

func (r *HumanAgentRepository) Delete(id uint) error {
	return r.db.Delete(&models.HumanAgent{}, id).Error
}

func (r *HumanAgentRepository) RevokeAccess(id uint) error {
	return r.db.Model(&models.HumanAgent{}).Where("id = ?", id).Update("active", false).Error
}
