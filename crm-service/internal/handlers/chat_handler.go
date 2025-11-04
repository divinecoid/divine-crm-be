package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/services"
	"divine-crm/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ChatHandler handles chat HTTP requests
type ChatHandler struct {
	service *services.ChatService
}

// NewChatHandler creates a new chat handler
func NewChatHandler(service *services.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

// GetAll handles GET /chats
func (h *ChatHandler) GetAll(c *fiber.Ctx) error {
	messages, err := h.service.GetAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch chats", err)
	}
	return utils.SuccessResponse(c, messages)
}

// GetByID handles GET /chats/:id
func (h *ChatHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	message, err := h.service.GetByID(uint(id))
	if err != nil {
		return utils.NotFoundResponse(c, "Chat")
	}
	return utils.SuccessResponse(c, message)
}

// Create handles POST /chats
func (h *ChatHandler) Create(c *fiber.Ctx) error {
	var message models.ChatMessage
	if err := c.BodyParser(&message); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.service.Create(&message); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create chat", err)
	}

	return c.Status(fiber.StatusCreated).JSON(utils.StandardResponse{
		Success: true,
		Data:    message,
	})
}

// Update handles PUT /chats/:id
func (h *ChatHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	message, err := h.service.GetByID(uint(id))
	if err != nil {
		return utils.NotFoundResponse(c, "Chat")
	}

	if err := c.BodyParser(message); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.service.Update(message); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update chat", err)
	}

	return utils.SuccessResponse(c, message)
}

// Assign handles POST /chats/:id/assign
func (h *ChatHandler) Assign(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	var req struct {
		AssignedTo    string `json:"assigned_to"`
		AssignedAgent string `json:"assigned_agent"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.service.Assign(uint(id), req.AssignedTo, req.AssignedAgent); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to assign chat", err)
	}

	message, _ := h.service.GetByID(uint(id))
	return utils.SuccessResponse(c, message)
}

// Resolve handles POST /chats/:id/resolve
func (h *ChatHandler) Resolve(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	if err := h.service.Resolve(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to resolve chat", err)
	}

	message, _ := h.service.GetByID(uint(id))
	return utils.SuccessResponse(c, message)
}
