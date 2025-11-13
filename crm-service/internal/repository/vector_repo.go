package repository

import (
	"divine-crm/internal/models"
	"fmt"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type VectorRepository struct {
	db *gorm.DB
}

func NewVectorRepository(db *gorm.DB) *VectorRepository {
	return &VectorRepository{db: db}
}

// ==================== KNOWLEDGE BASE ====================

func (r *VectorRepository) CreateKnowledge(kb *models.KnowledgeBase) error {
	return r.db.Create(kb).Error
}

func (r *VectorRepository) GetAllKnowledge() ([]models.KnowledgeBase, error) {
	var knowledge []models.KnowledgeBase
	err := r.db.Where("active = ?", true).Find(&knowledge).Error
	return knowledge, err
}

func (r *VectorRepository) FindSimilarKnowledge(embedding pgvector.Vector, limit int) ([]models.KnowledgeBase, error) {
	var results []models.KnowledgeBase

	query := `
		SELECT *, 
		       embedding <=> ? as distance
		FROM knowledge_bases 
		WHERE active = true
		ORDER BY embedding <=> ?
		LIMIT ?
	`

	err := r.db.Raw(query, embedding, embedding, limit).Scan(&results).Error
	return results, err
}

// ==================== CHAT HISTORY ====================

func (r *VectorRepository) SaveChatHistory(history *models.ChatHistory) error {
	return r.db.Create(history).Error
}

func (r *VectorRepository) FindSimilarChats(contactID uint, embedding pgvector.Vector, limit int) ([]models.ChatHistory, error) {
	var results []models.ChatHistory

	query := `
		SELECT *,
		       message_embedding <=> ? as distance
		FROM chat_histories
		WHERE contact_id = ?
		ORDER BY message_embedding <=> ?
		LIMIT ?
	`

	err := r.db.Raw(query, embedding, contactID, embedding, limit).Scan(&results).Error
	return results, err
}

// ==================== PRODUCT EMBEDDINGS ====================

func (r *VectorRepository) CreateProductEmbedding(pe *models.ProductEmbedding) error {
	return r.db.Create(pe).Error
}

func (r *VectorRepository) FindSimilarProducts(embedding pgvector.Vector, limit int) ([]models.ProductEmbedding, error) {
	var results []models.ProductEmbedding

	query := `
		SELECT *,
		       embedding <=> ? as distance
		FROM product_embeddings
		ORDER BY embedding <=> ?
		LIMIT ?
	`

	err := r.db.Raw(query, embedding, embedding, limit).Scan(&results).Error
	return results, err
}

// ==================== FAQ ====================

func (r *VectorRepository) CreateFAQ(faq *models.FAQEmbedding) error {
	return r.db.Create(faq).Error
}

func (r *VectorRepository) GetAllFAQ() ([]models.FAQEmbedding, error) {
	var faqs []models.FAQEmbedding
	err := r.db.Where("active = ?", true).Order("hit_count DESC").Find(&faqs).Error
	return faqs, err
}

func (r *VectorRepository) FindSimilarFAQ(embedding pgvector.Vector, limit int) ([]models.FAQEmbedding, error) {
	var results []models.FAQEmbedding

	query := `
		SELECT *,
		       embedding <=> ? as distance
		FROM faq_embeddings
		WHERE active = true
		ORDER BY embedding <=> ?
		LIMIT ?
	`

	err := r.db.Raw(query, embedding, embedding, limit).Scan(&results).Error
	return results, err
}

func (r *VectorRepository) IncrementFAQHit(id uint) error {
	return r.db.Model(&models.FAQEmbedding{}).
		Where("id = ?", id).
		Update("hit_count", gorm.Expr("hit_count + 1")).Error
}

// ==================== SIMPLE KEYWORD SEARCH (Fallback) ====================

func (r *VectorRepository) SearchKnowledgeByKeyword(query string, limit int) ([]models.KnowledgeBase, error) {
	var results []models.KnowledgeBase

	searchTerm := fmt.Sprintf("%%%s%%", query)
	err := r.db.Where("active = ?", true).
		Where("title ILIKE ? OR content ILIKE ? OR tags ILIKE ?",
			searchTerm, searchTerm, searchTerm).
		Limit(limit).
		Find(&results).Error

	return results, err
}
