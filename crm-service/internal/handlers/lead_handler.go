package handlers

import (
	"divine-crm/internal/services"
	"github.com/gofiber/fiber/v2"
)

// LeadHandler is DEPRECATED
// Use ContactHandler instead with filtering by ContactStatus = "Leads"
type LeadHandler struct {
	contactHandler *ContactHandler
}

// NewLeadHandler creates a lead handler (deprecated - use ContactHandler)
func NewLeadHandler(contactService *services.ContactService) *LeadHandler {
	return &LeadHandler{
		contactHandler: NewContactHandler(contactService),
	}
}

// GetAll returns all leads (contacts with status = "Leads")
func (h *LeadHandler) GetAll(c *fiber.Ctx) error {
	// Set query parameter to filter by Leads
	c.Request().URI().SetQueryString("status=Leads")
	return h.contactHandler.GetByStatus(c)
}

// GetByID returns a lead by ID
func (h *LeadHandler) GetByID(c *fiber.Ctx) error {
	return h.contactHandler.GetByID(c)
}

// Create creates a new lead
func (h *LeadHandler) Create(c *fiber.Ctx) error {
	// This would need special handling to set ContactStatus = "Leads"
	return c.Status(fiber.StatusGone).JSON(fiber.Map{
		"error": "This endpoint is deprecated. Use /api/v1/contacts instead with contact_status=Leads",
	})
}

// Update updates a lead
func (h *LeadHandler) Update(c *fiber.Ctx) error {
	return h.contactHandler.Update(c)
}

// Delete deletes a lead
func (h *LeadHandler) Delete(c *fiber.Ctx) error {
	return h.contactHandler.Delete(c)
}
