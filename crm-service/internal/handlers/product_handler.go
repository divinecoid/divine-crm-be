package handlers

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	repo *repository.ProductRepository
}

// NewProductHandler creates a new product handler
func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// GetAll handles GET /products
func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	products, err := h.repo.FindAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch products", err)
	}
	return utils.SuccessResponse(c, products)
}

// GetByID handles GET /products/:id
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	product, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.NotFoundResponse(c, "Product")
	}
	return utils.SuccessResponse(c, product)
}

// Create handles POST /products
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.repo.Create(&product); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create product", err)
	}

	return c.Status(fiber.StatusCreated).JSON(utils.StandardResponse{
		Success: true,
		Data:    product,
	})
}

// Update handles PUT /products/:id
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	product, err := h.repo.FindByID(uint(id))
	if err != nil {
		return utils.NotFoundResponse(c, "Product")
	}

	if err := c.BodyParser(product); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.repo.Update(product); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update product", err)
	}

	return utils.SuccessResponse(c, product)
}

// Delete handles DELETE /products/:id
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete product", err)
	}

	return utils.SuccessMessageResponse(c, "Product deleted successfully", nil)
}
