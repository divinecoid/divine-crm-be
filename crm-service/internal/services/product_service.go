package services

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ProductService struct {
	repo   *repository.ProductRepository
	logger *utils.Logger
}

func NewProductService(repo *repository.ProductRepository, logger *utils.Logger) *ProductService {
	return &ProductService{
		repo:   repo,
		logger: logger,
	}
}

// GetAll returns all products
func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.FindAll()
}

// GetByID returns a product by ID
func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	return s.repo.FindByID(id)
}

// GetActiveProducts returns all active products with stock
func (s *ProductService) GetActiveProducts() ([]models.Product, error) {
	return s.repo.FindActive()
}

// SearchProducts searches products by name or code
func (s *ProductService) SearchProducts(query string) ([]models.Product, error) {
	return s.repo.Search(query)
}

// Create creates a new product
func (s *ProductService) Create(product *models.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	s.logger.Info("Creating product", "code", product.Code, "name", product.Name)
	return s.repo.Create(product)
}

// Update updates a product
func (s *ProductService) Update(product *models.Product) error {
	product.UpdatedAt = time.Now()

	s.logger.Info("Updating product", "id", product.ID, "name", product.Name)
	return s.repo.Update(product)
}

// Delete deletes a product
func (s *ProductService) Delete(id uint) error {
	s.logger.Info("Deleting product", "id", id)
	return s.repo.Delete(id)
}

// FormatProductList formats products for AI context
func (s *ProductService) FormatProductList() string {
	products, err := s.GetActiveProducts()
	if err != nil || len(products) == 0 {
		return "Saat ini tidak ada produk yang tersedia."
	}

	var builder strings.Builder
	builder.WriteString("\n=== DAFTAR PRODUK YANG TERSEDIA ===\n\n")

	for i, product := range products {
		builder.WriteString(fmt.Sprintf("%d. %s (Kode: %s)\n", i+1, product.Name, product.Code))
		builder.WriteString(fmt.Sprintf("   💰 Harga: Rp %s\n", formatCurrency(product.Price)))
		builder.WriteString(fmt.Sprintf("   📦 Stok: %d unit\n", product.Stock))
		if product.Description != "" {
			builder.WriteString(fmt.Sprintf("   📝 Deskripsi: %s\n", product.Description))
		}
		builder.WriteString("\n")
	}

	builder.WriteString("=== CARA ORDER ===\n")
	builder.WriteString("Customer dapat order dengan menyebutkan nama atau kode produk.\n")
	builder.WriteString("Contoh: 'Saya mau pesan Produk 1' atau 'Order F003'\n")

	return builder.String()
}

// Helper function to format currency
func formatCurrency(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 0, 64)
}
