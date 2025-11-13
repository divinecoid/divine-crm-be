package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type BroadcastHandler struct {
	service *services.BroadcastService
}

func NewBroadcastHandler(service *services.BroadcastService) *BroadcastHandler {
	return &BroadcastHandler{service: service}
}

// Templates
func (h *BroadcastHandler) GetAllTemplates(c *fiber.Ctx) error {
	templates, err := h.service.GetAllTemplates()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": templates})
}

func (h *BroadcastHandler) GetTemplateByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	template, err := h.service.GetTemplateByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Template not found"})
	}

	return c.JSON(fiber.Map{"data": template})
}

func (h *BroadcastHandler) CreateTemplate(c *fiber.Ctx) error {
	var template models.BroadcastTemplate
	if err := c.BodyParser(&template); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.CreateTemplate(&template); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": template})
}

func (h *BroadcastHandler) UpdateTemplate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var template models.BroadcastTemplate
	if err := c.BodyParser(&template); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	template.ID = uint(id)
	if err := h.service.UpdateTemplate(&template); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": template})
}

func (h *BroadcastHandler) DeleteTemplate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.DeleteTemplate(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Template deleted successfully"})
}

// Broadcasting
func (h *BroadcastHandler) SendBroadcast(c *fiber.Ctx) error {
	var body struct {
		TemplateID uint   `json:"template_id"`
		SentBy     string `json:"sent_by"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.SendBroadcast(body.TemplateID, body.SentBy); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Broadcast started successfully"})
}

// History
func (h *BroadcastHandler) GetHistory(c *fiber.Ctx) error {
	history, err := h.service.GetHistory()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": history})
}
