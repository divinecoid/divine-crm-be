package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check handles health check requests
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	// Check database connection
	sqlDB, err := h.db.DB()
	dbStatus := "healthy"
	if err != nil {
		dbStatus = "unhealthy"
	} else {
		if err := sqlDB.Ping(); err != nil {
			dbStatus = "unhealthy"
		}
	}

	return c.JSON(fiber.Map{
		"status":    "ok",
		"service":   "crm-service",
		"database":  dbStatus,
		"timestamp": time.Now(),
	})
}
