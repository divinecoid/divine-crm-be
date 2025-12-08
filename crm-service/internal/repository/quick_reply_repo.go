package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
)

type QuickReplyRepository struct {
	db *gorm.DB
}

func NewQuickReplyRepository(db *gorm.DB) *QuickReplyRepository {
	return &QuickReplyRepository{db: db}
}

func (r *QuickReplyRepository) FindAll() ([]models.QuickReply, error) {
	var replies []models.QuickReply
	err := r.db.Where("active = ?", true).Order("trigger ASC").Find(&replies).Error
	return replies, err
}

func (r *QuickReplyRepository) FindByTrigger(trigger string) (*models.QuickReply, error) {
	var reply models.QuickReply
	err := r.db.Where("trigger = ? AND active = ?", trigger, true).First(&reply).Error
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func (r *QuickReplyRepository) Create(reply *models.QuickReply) error {
	return r.db.Create(reply).Error
}

func (r *QuickReplyRepository) Update(reply *models.QuickReply) error {
	return r.db.Save(reply).Error
}

func (r *QuickReplyRepository) Delete(id uint) error {
	return r.db.Delete(&models.QuickReply{}, id).Error
}
