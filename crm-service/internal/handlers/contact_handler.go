package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/services"
	"divine-crm/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ContactHandler handles contact HTTP requests
type ContactHandler struct {
	service *services.ContactService
}

// NewContactHandler creates a new contact handler
func NewContactHandler(service *services.ContactService) *ContactHandler {
	return &ContactHandler{service: service}
}

// GetAll handles GET /contacts
func (h *ContactHandler) GetAll(c *fiber.Ctx) error {
	contacts, err := h.service.GetAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch contacts", err)
	}
	return utils.SuccessResponse(c, contacts)
}

// GetByID handles GET /contacts/:id
func (h *ContactHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	contact, err := h.service.GetByID(uint(id))
	if err != nil {
		return utils.NotFoundResponse(c, "Contact")
	}
	return utils.SuccessResponse(c, contact)
}

// Create handles POST /contacts
func (h *ContactHandler) Create(c *fiber.Ctx) error {
	var contact models.Contact
	if err := c.BodyParser(&contact); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.service.Create(&contact); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create contact", err)
	}

	return c.Status(fiber.StatusCreated).JSON(utils.StandardResponse{
		Success: true,
		Data:    contact,
	})
}

// Update handles PUT /contacts/:id
func (h *ContactHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	contact, err := h.service.GetByID(uint(id))
	if err != nil {
		return utils.NotFoundResponse(c, "Contact")
	}

	if err := c.BodyParser(contact); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.service.Update(contact); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update contact", err)
	}

	return utils.SuccessResponse(c, contact)
}

// Delete handles DELETE /contacts/:id
func (h *ContactHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete contact", err)
	}

	return utils.SuccessMessageResponse(c, "Contact deleted successfully", nil)
}
