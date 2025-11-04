package handlers

import (
	"divine-crm/internal/repository"
	"divine-crm/internal/services"
	"divine-crm/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AIHandler handles AI HTTP requests
type AIHandler struct {
	service *services.AIService
	repo    *repository.AIRepository
}

// NewAIHandler creates a new AI handler
func NewAIHandler(service *services.AIService, repo *repository.AIRepository) *AIHandler {
	return &AIHandler{
		service: service,
		repo:    repo,
	}
}

// Process handles POST /ai/process
func (h *AIHandler) Process(c *fiber.Ctx) error {
	var req struct {
		Message string `json:"message"`
		AgentID uint   `json:"agent_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	response, tokensUsed, err := h.service.ProcessMessage(req.Message, req.AgentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "AI processing failed", err)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"response":    response,
		"tokens_used": tokensUsed,
	})
}

// GetTokenBalances handles GET /ai/token-balance
func (h *AIHandler) GetTokenBalances(c *fiber.Ctx) error {
	balances, err := h.repo.GetTokenBalances()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch token balances", err)
	}
	return utils.SuccessResponse(c, balances)
}

// GetConfigurations handles GET /ai/configurations
func (h *AIHandler) GetConfigurations(c *fiber.Ctx) error {
	configs, err := h.repo.GetAllConfigurations()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch AI configurations", err)
	}
	return utils.SuccessResponse(c, configs)
}

// GetAgents handles GET /ai/agents
func (h *AIHandler) GetAgents(c *fiber.Ctx) error {
	agents, err := h.repo.GetAllAgents()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch AI agents", err)
	}
	return utils.SuccessResponse(c, agents)
}
