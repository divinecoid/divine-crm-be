package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type QuickReplyHandler struct {
	service *services.QuickReplyService
}

func NewQuickReplyHandler(service *services.QuickReplyService) *QuickReplyHandler {
	return &QuickReplyHandler{service: service}
}

func (h *QuickReplyHandler) GetAll(c *fiber.Ctx) error {
	replies, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": replies})
}

func (h *QuickReplyHandler) Create(c *fiber.Ctx) error {
	var reply models.QuickReply
	if err := c.BodyParser(&reply); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Create(&reply); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": reply})
}

func (h *QuickReplyHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var reply models.QuickReply
	if err := c.BodyParser(&reply); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	reply.ID = uint(id)
	if err := h.service.Update(&reply); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": reply})
}

func (h *QuickReplyHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Quick reply deleted successfully"})
}
