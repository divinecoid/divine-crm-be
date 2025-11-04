package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// LeadHandler handles lead HTTP requests
type LeadHandler struct {
	repo *repository.LeadRepository
}

// NewLeadHandler creates a new lead handler
func NewLeadHandler(repo *repository.LeadRepository) *LeadHandler {
	return &LeadHandler{repo: repo}
}

// GetAll handles GET /leads
func (h *LeadHandler) GetAll(c *fiber.Ctx) error {
	leads, err := h.repo.FindAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch leads", err)
	}
	return utils.SuccessResponse(c, leads)
}

// GetByID handles GET /leads/:id
func (h *LeadHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	lead, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.NotFoundResponse(c, "Lead")
	}
	return utils.SuccessResponse(c, lead)
}

// Create handles POST /leads
func (h *LeadHandler) Create(c *fiber.Ctx) error {
	var lead models.Lead
	if err := c.BodyParser(&lead); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.repo.Create(&lead); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create lead", err)
	}

	return c.Status(fiber.StatusCreated).JSON(utils.StandardResponse{
		Success: true,
		Data:    lead,
	})
}

// Update handles PUT /leads/:id
func (h *LeadHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	lead, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.NotFoundResponse(c, "Lead")
	}

	if err := c.BodyParser(lead); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.repo.Update(lead); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update lead", err)
	}

	return utils.SuccessResponse(c, lead)
}

// Delete handles DELETE /leads/:id
func (h *LeadHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete lead", err)
	}

	return utils.SuccessMessageResponse(c, "Lead deleted successfully", nil)
}
