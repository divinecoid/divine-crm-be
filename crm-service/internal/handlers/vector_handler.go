package handlers

import (
	"divine-crm/internal/services"
	"divine-crm/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type VectorHandler struct {
	service *services.VectorService
}

func NewVectorHandler(service *services.VectorService) *VectorHandler {
	return &VectorHandler{service: service}
}

// ==================== KNOWLEDGE BASE ====================

func (h *VectorHandler) AddKnowledge(c *fiber.Ctx) error {
	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
		Tags     string `json:"tags"`
		Source   string `json:"source"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if err := h.service.AddKnowledge(req.Title, req.Content, req.Category, req.Tags, req.Source); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Knowledge added successfully"})
}

func (h *VectorHandler) SearchKnowledge(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return utils.BadRequestResponse(c, "Query parameter 'q' is required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "5"))

	results, err := h.service.SearchKnowledge(query, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, results)
}

func (h *VectorHandler) GetAllKnowledge(c *fiber.Ctx) error {
	results, err := h.service.GetAllKnowledge()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, results)
}

// ==================== PRODUCT SEARCH ====================

func (h *VectorHandler) SearchProducts(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return utils.BadRequestResponse(c, "Query parameter 'q' is required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "5"))

	results, err := h.service.SearchProducts(query, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, results)
}

func (h *VectorHandler) AddProductEmbedding(c *fiber.Ctx) error {
	var req struct {
		ProductID   uint   `json:"product_id"`
		Description string `json:"description"`
		Features    string `json:"features"`
		UseCases    string `json:"use_cases"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if err := h.service.AddProductEmbedding(req.ProductID, req.Description, req.Features, req.UseCases); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Product embedding added successfully"})
}

// ==================== FAQ ====================

func (h *VectorHandler) AddFAQ(c *fiber.Ctx) error {
	var req struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
		Category string `json:"category"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	if err := h.service.AddFAQ(req.Question, req.Answer, req.Category); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "FAQ added successfully"})
}

func (h *VectorHandler) SearchFAQ(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return utils.BadRequestResponse(c, "Query parameter 'q' is required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "3"))

	results, err := h.service.SearchFAQ(query, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, results)
}

func (h *VectorHandler) GetAllFAQ(c *fiber.Ctx) error {
	results, err := h.service.GetAllFAQ()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, results)
}

// ==================== CHAT CONTEXT ====================

func (h *VectorHandler) GetSimilarConversations(c *fiber.Ctx) error {
	contactID, _ := strconv.Atoi(c.Params("contactId"))
	query := c.Query("q")
	limit, _ := strconv.Atoi(c.Query("limit", "5"))

	results, err := h.service.GetSimilarConversations(uint(contactID), query, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, results)
}
