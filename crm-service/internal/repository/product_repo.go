package repository

import (
	"divine-crm/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Order("name ASC").Find(&products).Error
	return products, err
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *ProductRepository) FindActive() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("stock > ?", 0).Order("name ASC").Find(&products).Error
	return products, err
}

func (r *ProductRepository) Search(query string) ([]models.Product, error) {
	var products []models.Product
	searchPattern := "%" + query + "%"
	err := r.db.Where("name ILIKE ? OR code ILIKE ?", searchPattern, searchPattern).
		Where("stock > ?", 0).
		Find(&products).Error
	return products, err
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
