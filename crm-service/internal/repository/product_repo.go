package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
)

// ProductRepository handles product data operations
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// FindAll returns all products
func (r *ProductRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Order("created_at desc").Find(&products).Error
	return products, err
}

// FindByID returns a product by ID
func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Create creates a new product
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// Update updates a product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete deletes a product
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
