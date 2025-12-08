package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type AIAgentHandler struct {
	service *services.AIAgentService
}

func NewAIAgentHandler(service *services.AIAgentService) *AIAgentHandler {
	return &AIAgentHandler{service: service}
}

func (h *AIAgentHandler) GetAll(c *fiber.Ctx) error {
	agents, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": agents})
}

func (h *AIAgentHandler) GetActive(c *fiber.Ctx) error {
	agent, err := h.service.GetActive()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": agent})
}

func (h *AIAgentHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	agent, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Agent not found"})
	}

	return c.JSON(fiber.Map{"data": agent})
}

func (h *AIAgentHandler) Create(c *fiber.Ctx) error {
	var agent models.AIAgent
	if err := c.BodyParser(&agent); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Create(&agent); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": agent})
}

func (h *AIAgentHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var agent models.AIAgent
	if err := c.BodyParser(&agent); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	agent.ID = uint(id)
	if err := h.service.Update(&agent); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": agent})
}

func (h *AIAgentHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Agent deleted successfully"})
}

func (h *AIAgentHandler) Activate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.Activate(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Agent activated successfully"})
}
