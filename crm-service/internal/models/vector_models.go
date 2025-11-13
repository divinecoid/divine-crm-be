package models

import (
	"github.com/pgvector/pgvector-go"
	"time"
)

// KnowledgeBase stores company knowledge with vector embeddings
type KnowledgeBase struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Title     string          `gorm:"size:255;not null" json:"title"`
	Content   string          `gorm:"type:text;not null" json:"content"`
	Category  string          `gorm:"size:100" json:"category"`
	Tags      string          `gorm:"size:255" json:"tags"`
	Source    string          `gorm:"size:255" json:"source"`
	Embedding pgvector.Vector `gorm:"type:vector(1536)" json:"-"`
	Active    bool            `gorm:"default:true" json:"active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// ChatHistory stores chat messages with embeddings for RAG
type ChatHistory struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	ContactID         uint            `gorm:"not null;index" json:"contact_id"`
	UserMessage       string          `gorm:"type:text;not null" json:"user_message"`
	AIResponse        string          `gorm:"type:text" json:"ai_response"`
	MessageEmbedding  pgvector.Vector `gorm:"type:vector(1536)" json:"-"`
	ResponseEmbedding pgvector.Vector `gorm:"type:vector(1536)" json:"-"`
	Sentiment         string          `gorm:"size:50" json:"sentiment"` // positive, negative, neutral
	Intent            string          `gorm:"size:100" json:"intent"`   // inquiry, complaint, order, etc
	CreatedAt         time.Time       `json:"created_at"`
}

// ProductEmbedding stores product info with embeddings for semantic search
type ProductEmbedding struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	ProductID   uint            `gorm:"not null;index" json:"product_id"`
	Description string          `gorm:"type:text" json:"description"`
	Features    string          `gorm:"type:text" json:"features"`
	UseCases    string          `gorm:"type:text" json:"use_cases"`
	Embedding   pgvector.Vector `gorm:"type:vector(1536)" json:"-"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// FAQEmbedding stores FAQs with embeddings
type FAQEmbedding struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Question  string          `gorm:"type:text;not null" json:"question"`
	Answer    string          `gorm:"type:text;not null" json:"answer"`
	Category  string          `gorm:"size:100" json:"category"`
	Embedding pgvector.Vector `gorm:"type:vector(1536)" json:"-"`
	HitCount  int             `gorm:"default:0" json:"hit_count"` // Track usage
	Active    bool            `gorm:"default:true" json:"active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
