package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type HumanAgentHandler struct {
	service *services.HumanAgentService
}

func NewHumanAgentHandler(service *services.HumanAgentService) *HumanAgentHandler {
	return &HumanAgentHandler{service: service}
}

func (h *HumanAgentHandler) GetAll(c *fiber.Ctx) error {
	agents, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": agents})
}

func (h *HumanAgentHandler) GetActive(c *fiber.Ctx) error {
	agents, err := h.service.GetActive()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": agents})
}

func (h *HumanAgentHandler) GetByID(c *fiber.Ctx) error {
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

func (h *HumanAgentHandler) Create(c *fiber.Ctx) error {
	var agent models.HumanAgent
	if err := c.BodyParser(&agent); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Create(&agent); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": agent})
}

func (h *HumanAgentHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var agent models.HumanAgent
	if err := c.BodyParser(&agent); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	agent.ID = uint(id)
	if err := h.service.Update(&agent); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": agent})
}

func (h *HumanAgentHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Agent deleted successfully"})
}

func (h *HumanAgentHandler) RevokeAccess(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.RevokeAccess(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Access revoked successfully"})
}

func (h *HumanAgentHandler) ResetPassword(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var body struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.ResetPassword(uint(id), body.NewPassword); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Password reset successfully"})
}

func (h *HumanAgentHandler) Login(c *fiber.Ctx) error {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	agent, err := h.service.Authenticate(credentials.Username, credentials.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate JWT token
	token, err := h.service.GenerateToken(agent)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": fiber.Map{
			"user":  agent,
			"token": token,
		},
	})
}
