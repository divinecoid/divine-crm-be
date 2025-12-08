package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ChatHandler struct {
	service *services.ChatService
}

func NewChatHandler(service *services.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

// GetAll returns all chat messages
func (h *ChatHandler) GetAll(c *fiber.Ctx) error {
	messages, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": messages})
}

// GetByID returns a chat message by ID
func (h *ChatHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	message, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Message not found"})
	}

	return c.JSON(fiber.Map{"data": message})
}

// GetUnassigned returns all unassigned chat messages
func (h *ChatHandler) GetUnassigned(c *fiber.Ctx) error {
	messages, err := h.service.GetByStatus("Unassigned")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"data":  messages,
		"count": len(messages),
	})
}

// GetAssigned returns all assigned chat messages
func (h *ChatHandler) GetAssigned(c *fiber.Ctx) error {
	messages, err := h.service.GetByStatus("Assigned")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"data":  messages,
		"count": len(messages),
	})
}

// GetResolved returns all resolved chat messages
func (h *ChatHandler) GetResolved(c *fiber.Ctx) error {
	messages, err := h.service.GetByStatus("Resolved")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"data":  messages,
		"count": len(messages),
	})
}

// Create creates a new chat message
func (h *ChatHandler) Create(c *fiber.Ctx) error {
	var message models.ChatMessage
	if err := c.BodyParser(&message); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Create(&message); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": message})
}

// Update updates a chat message
func (h *ChatHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var message models.ChatMessage
	if err := c.BodyParser(&message); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	message.ID = uint(id)
	if err := h.service.Update(&message); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": message})
}

// Assign assigns a chat to an agent
func (h *ChatHandler) Assign(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var body struct {
		AssignedTo    string `json:"assigned_to"`
		AssignedAgent string `json:"assigned_agent"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Assign(uint(id), body.AssignedTo, body.AssignedAgent); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Chat assigned successfully"})
}

// Resolve marks a chat as resolved
func (h *ChatHandler) Resolve(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.Resolve(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Chat resolved successfully"})
}

// TakeOver - Human agent takes over from AI
func (h *ChatHandler) TakeOver(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var body struct {
		AgentName string `json:"agent_name"` // Human agent name
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.TakeOver(uint(id), body.AgentName); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Chat taken over by human agent successfully",
		"agent":   body.AgentName,
	})
}

// BackToAI - Return chat back to AI
func (h *ChatHandler) BackToAI(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := h.service.BackToAI(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Chat returned to AI successfully",
	})
}

// GetByContact returns all chat messages for a specific contact
func (h *ChatHandler) GetByContact(c *fiber.Ctx) error {
	contactID, err := strconv.Atoi(c.Query("contact_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid contact_id"})
	}

	messages, err := h.service.GetByContactID(uint(contactID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  messages,
		"count": len(messages),
	})
}

// GetByChannel returns chat messages filtered by channel
func (h *ChatHandler) GetByChannel(c *fiber.Ctx) error {
	channel := c.Query("channel")
	if channel == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Channel parameter required"})
	}

	messages, err := h.service.GetByChannel(channel)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":    messages,
		"count":   len(messages),
		"channel": channel,
	})
}

// GetStats returns chat statistics
func (h *ChatHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.service.GetStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": stats})
}

// AddLabel adds a label to a chat message
func (h *ChatHandler) AddLabel(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var body struct {
		LabelID string `json:"label_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.AddLabel(uint(id), body.LabelID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Label added successfully"})
}
