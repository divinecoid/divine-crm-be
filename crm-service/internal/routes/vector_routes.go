package routes

import (
	"divine-crm/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupVectorRoutes(api fiber.Router, h *handlers.VectorHandler) {
	vectors := api.Group("/vectors")

	// Knowledge Base
	vectors.Post("/knowledge", h.AddKnowledge)
	vectors.Get("/knowledge", h.GetAllKnowledge)
	vectors.Get("/knowledge/search", h.SearchKnowledge)

	// FAQ
	vectors.Post("/faq", h.AddFAQ)
	vectors.Get("/faq", h.GetAllFAQ)
	vectors.Get("/faq/search", h.SearchFAQ)

	// Product Semantic Search
	vectors.Post("/products/embedding", h.AddProductEmbedding)
	vectors.Get("/products/search", h.SearchProducts)

	// Chat History
	vectors.Get("/chat-history/:contactId", h.GetSimilarConversations)
}
