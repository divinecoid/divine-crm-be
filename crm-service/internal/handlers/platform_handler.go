package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type PlatformHandler struct {
	service *services.PlatformService
}

func NewPlatformHandler(service *services.PlatformService) *PlatformHandler {
	return &PlatformHandler{service: service}
}

func (h *PlatformHandler) GetAll(c *fiber.Ctx) error {
	platforms, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": platforms})
}

func (h *PlatformHandler) GetActive(c *fiber.Ctx) error {
	platforms, err := h.service.GetActive()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": platforms})
}

func (h *PlatformHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	platform, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Platform not found"})
	}

	return c.JSON(fiber.Map{"data": platform})
}

func (h *PlatformHandler) Create(c *fiber.Ctx) error {
	var platform models.ConnectedPlatform
	if err := c.BodyParser(&platform); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Create(&platform); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": platform})
}

func (h *PlatformHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var platform models.ConnectedPlatform
	if err := c.BodyParser(&platform); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	platform.ID = uint(id)
	if err := h.service.Update(&platform); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": platform})
}

func (h *PlatformHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Platform deleted successfully"})
}
