package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ChatLabelHandler struct {
	service *services.ChatLabelService
}

func NewChatLabelHandler(service *services.ChatLabelService) *ChatLabelHandler {
	return &ChatLabelHandler{service: service}
}

func (h *ChatLabelHandler) GetAll(c *fiber.Ctx) error {
	labels, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": labels})
}

func (h *ChatLabelHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	label, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Label not found"})
	}

	return c.JSON(fiber.Map{"data": label})
}

func (h *ChatLabelHandler) Create(c *fiber.Ctx) error {
	var label models.ChatLabel
	if err := c.BodyParser(&label); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Create(&label); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": label})
}

func (h *ChatLabelHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var label models.ChatLabel
	if err := c.BodyParser(&label); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	label.ID = uint(id)
	if err := h.service.Update(&label); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": label})
}

func (h *ChatLabelHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Label deleted successfully"})
}
