package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
	"time"
)

type AnalyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) FindByDateRange(start, end time.Time) ([]models.Analytics, error) {
	var analytics []models.Analytics
	err := r.db.Where("date BETWEEN ? AND ?", start, end).Order("date DESC").Find(&analytics).Error
	return analytics, err
}

func (r *AnalyticsRepository) FindByDate(date time.Time) (*models.Analytics, error) {
	var analytic models.Analytics
	err := r.db.Where("date = ?", date.Format("2006-01-02")).First(&analytic).Error
	if err != nil {
		return nil, err
	}
	return &analytic, nil
}

func (r *AnalyticsRepository) CreateOrUpdate(analytic *models.Analytics) error {
	var existing models.Analytics
	err := r.db.Where("date = ?", analytic.Date.Format("2006-01-02")).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(analytic).Error
	}

	// Update existing
	existing.TotalMessages = analytic.TotalMessages
	existing.UnassignedChats = analytic.UnassignedChats
	existing.AssignedChats = analytic.AssignedChats
	existing.ResolvedChats = analytic.ResolvedChats
	existing.NewContacts = analytic.NewContacts
	existing.TotalTokensUsed = analytic.TotalTokensUsed
	existing.WhatsAppMessages = analytic.WhatsAppMessages
	existing.InstagramMessages = analytic.InstagramMessages
	existing.TelegramMessages = analytic.TelegramMessages

	return r.db.Save(&existing).Error
}

func (r *AnalyticsRepository) GetTodayStats() (*models.Analytics, error) {
	return r.FindByDate(time.Now())
}
