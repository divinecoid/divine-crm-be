package routes

import (
	"divine-crm/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupContactRoutes(api fiber.Router, h *handlers.ContactHandler) {
	contacts := api.Group("/contacts")
	contacts.Get("/", h.GetAll)
	contacts.Get("/stats", h.GetStats)
	contacts.Get("/search", h.Search)
	contacts.Get("/status", h.GetByStatus)
	contacts.Get("/temperature", h.GetByTemperature)
	contacts.Get("/:id", h.GetByID)
	contacts.Post("/", h.Create)
	contacts.Put("/:id", h.Update)
	contacts.Delete("/:id", h.Delete)
	contacts.Patch("/:id/temperature", h.UpdateTemperature)
	contacts.Patch("/:id/status", h.UpdateStatus)
}
