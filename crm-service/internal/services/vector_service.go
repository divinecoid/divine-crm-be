package services

import (
	"bytes"
	"divine-crm/internal/config"
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pgvector/pgvector-go"
)

type VectorService struct {
	repo   *repository.VectorRepository
	config *config.Config
	logger *utils.Logger
}

func NewVectorService(repo *repository.VectorRepository, cfg *config.Config, logger *utils.Logger) *VectorService {
	return &VectorService{
		repo:   repo,
		config: cfg,
		logger: logger,
	}
}

// GenerateEmbedding generates vector embedding for text using OpenAI
func (s *VectorService) GenerateEmbedding(text string) (pgvector.Vector, error) {
	requestBody := map[string]interface{}{
		"input": text,
		"model": "text-embedding-ada-002",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return pgvector.Vector{}, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return pgvector.Vector{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.OpenAI.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return pgvector.Vector{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return pgvector.Vector{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return pgvector.Vector{}, fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return pgvector.Vector{}, err
	}

	if len(result.Data) == 0 {
		return pgvector.Vector{}, fmt.Errorf("no embedding returned")
	}

	return pgvector.NewVector(result.Data[0].Embedding), nil
}

// ==================== KNOWLEDGE BASE ====================

func (s *VectorService) AddKnowledge(title, content, category, tags, source string) error {
	embedding, err := s.GenerateEmbedding(content)
	if err != nil {
		return err
	}

	kb := &models.KnowledgeBase{
		Title:     title,
		Content:   content,
		Category:  category,
		Tags:      tags,
		Source:    source,
		Embedding: embedding,
		Active:    true,
	}

	return s.repo.CreateKnowledge(kb)
}

func (s *VectorService) SearchKnowledge(query string, limit int) ([]models.KnowledgeBase, error) {
	embedding, err := s.GenerateEmbedding(query)
	if err != nil {
		return nil, err
	}

	return s.repo.FindSimilarKnowledge(embedding, limit)
}

func (s *VectorService) GetAllKnowledge() ([]models.KnowledgeBase, error) {
	return s.repo.GetAllKnowledge()
}

// ==================== CHAT HISTORY ====================

func (s *VectorService) SaveChatWithEmbedding(contactID uint, userMessage, aiResponse, sentiment, intent string) error {
	msgEmbedding, err := s.GenerateEmbedding(userMessage)
	if err != nil {
		return err
	}

	respEmbedding, err := s.GenerateEmbedding(aiResponse)
	if err != nil {
		return err
	}

	history := &models.ChatHistory{
		ContactID:         contactID,
		UserMessage:       userMessage,
		AIResponse:        aiResponse,
		MessageEmbedding:  msgEmbedding,
		ResponseEmbedding: respEmbedding,
		Sentiment:         sentiment,
		Intent:            intent,
	}

	return s.repo.SaveChatHistory(history)
}

func (s *VectorService) GetSimilarConversations(contactID uint, query string, limit int) ([]models.ChatHistory, error) {
	embedding, err := s.GenerateEmbedding(query)
	if err != nil {
		return nil, err
	}

	return s.repo.FindSimilarChats(contactID, embedding, limit)
}

// ==================== PRODUCT EMBEDDINGS ====================

func (s *VectorService) AddProductEmbedding(productID uint, description, features, useCases string) error {
	combinedText := fmt.Sprintf("%s. %s. %s", description, features, useCases)
	embedding, err := s.GenerateEmbedding(combinedText)
	if err != nil {
		return err
	}

	pe := &models.ProductEmbedding{
		ProductID:   productID,
		Description: description,
		Features:    features,
		UseCases:    useCases,
		Embedding:   embedding,
	}

	return s.repo.CreateProductEmbedding(pe)
}

func (s *VectorService) SearchProducts(query string, limit int) ([]models.ProductEmbedding, error) {
	embedding, err := s.GenerateEmbedding(query)
	if err != nil {
		return nil, err
	}

	return s.repo.FindSimilarProducts(embedding, limit)
}

// ==================== FAQ ====================

func (s *VectorService) AddFAQ(question, answer, category string) error {
	embedding, err := s.GenerateEmbedding(question)
	if err != nil {
		return err
	}

	faq := &models.FAQEmbedding{
		Question:  question,
		Answer:    answer,
		Category:  category,
		Embedding: embedding,
		Active:    true,
	}

	return s.repo.CreateFAQ(faq)
}

func (s *VectorService) SearchFAQ(query string, limit int) ([]models.FAQEmbedding, error) {
	embedding, err := s.GenerateEmbedding(query)
	if err != nil {
		return nil, err
	}

	faqs, err := s.repo.FindSimilarFAQ(embedding, limit)
	if err != nil {
		return nil, err
	}

	// Increment hit count for best match
	if len(faqs) > 0 {
		s.repo.IncrementFAQHit(faqs[0].ID)
	}

	return faqs, nil
}

func (s *VectorService) GetAllFAQ() ([]models.FAQEmbedding, error) {
	return s.repo.GetAllFAQ()
}

// ==================== RAG (Retrieval Augmented Generation) ====================

func (s *VectorService) BuildRAGContext(query string) (string, error) {
	// Search knowledge base
	knowledge, err := s.SearchKnowledge(query, 3)
	if err != nil {
		s.logger.Warn("Failed to search knowledge", "error", err)
	}

	// Search FAQ
	faqs, err := s.SearchFAQ(query, 2)
	if err != nil {
		s.logger.Warn("Failed to search FAQ", "error", err)
	}

	// Build context
	var context string

	if len(knowledge) > 0 {
		context += "\n=== KNOWLEDGE BASE ===\n"
		for i, kb := range knowledge {
			context += fmt.Sprintf("%d. %s\n%s\n\n", i+1, kb.Title, kb.Content)
		}
	}

	if len(faqs) > 0 {
		context += "\n=== FAQ ===\n"
		for _, faq := range faqs {
			context += fmt.Sprintf("Q: %s\nA: %s\n\n", faq.Question, faq.Answer)
		}
	}

	return context, nil
}
