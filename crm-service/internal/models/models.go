package models

import (
	"time"
)

// Contact model
type Contact struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Code          string    `json:"code" gorm:"uniqueIndex;type:varchar(50)"`
	Channel       string    `json:"channel" gorm:"type:varchar(50)"`
	ChannelID     string    `json:"channel_id" gorm:"type:varchar(255);index"` // ✅ Simplified
	Name          string    `json:"name" gorm:"type:varchar(255)"`
	Temperature   string    `json:"temperature" gorm:"type:varchar(20);index"`
	FirstContact  time.Time `json:"first_contact"`
	LastContact   time.Time `json:"last_contact" gorm:"index"`
	LastAgent     string    `json:"last_agent" gorm:"type:varchar(255)"`
	LastAgentType string    `json:"last_agent_type" gorm:"type:varchar(50)"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Lead model
type Lead struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Code          string    `json:"code" gorm:"uniqueIndex;type:varchar(50)"`
	Channel       string    `json:"channel" gorm:"type:varchar(50)"`
	ChannelID     string    `json:"channel_id" gorm:"type:varchar(255)"`
	Temperature   string    `json:"temperature" gorm:"type:varchar(20)"`
	FirstContact  time.Time `json:"first_contact"`
	LastContact   time.Time `json:"last_contact"`
	LastAgent     string    `json:"last_agent" gorm:"type:varchar(255)"`
	LastAgentType string    `json:"last_agent_type" gorm:"type:varchar(50)"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Product model
type Product struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Code       string    `json:"code" gorm:"uniqueIndex;type:varchar(50)"`
	Name       string    `json:"name" gorm:"type:varchar(255)"`
	Price      float64   `json:"price"`
	Stock      int       `json:"stock"`
	UploadedBy string    `json:"uploaded_by" gorm:"type:varchar(255)"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ChatMessage model
type ChatMessage struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ContactID     uint      `json:"contact_id" gorm:"index"`
	ContactName   string    `json:"contact_name" gorm:"type:varchar(255)"`
	Message       string    `json:"message" gorm:"type:text"`
	Response      string    `json:"response" gorm:"type:text"`
	Status        string    `json:"status" gorm:"type:varchar(50);index;default:Unassigned"` // ✅ Fixed order
	AssignedTo    string    `json:"assigned_to" gorm:"type:varchar(255)"`
	AssignedAgent string    `json:"assigned_agent" gorm:"type:varchar(255)"`
	Channel       string    `json:"channel" gorm:"type:varchar(50)"`
	Labels        string    `json:"labels" gorm:"type:text"`
	TokensUsed    int       `json:"tokens_used" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at" gorm:"index"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ChatLabel model
type ChatLabel struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Label       string    `json:"label" gorm:"type:varchar(100)"`
	Description string    `json:"description" gorm:"type:varchar(500)"`
	Color       string    `json:"color" gorm:"type:varchar(50)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// QuickReply model
type QuickReply struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Trigger   string    `json:"trigger" gorm:"uniqueIndex;type:varchar(255)"`
	Response  string    `json:"response" gorm:"type:text"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BroadcastTemplate model
type BroadcastTemplate struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(255)"`
	Content   string    `json:"content" gorm:"type:text"`
	Type      string    `json:"type" gorm:"type:varchar(50)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AIConfiguration model
type AIConfiguration struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AIEngine  string    `json:"ai_engine" gorm:"uniqueIndex;type:varchar(100)"`
	Token     string    `json:"token" gorm:"type:text"`
	Endpoint  string    `json:"endpoint" gorm:"type:text"`
	Model     string    `json:"model" gorm:"type:varchar(100)"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AIAgent model
type AIAgent struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"uniqueIndex;type:varchar(255)"`
	AIEngine    string    `json:"ai_engine" gorm:"type:varchar(100)"`
	BasicPrompt string    `json:"basic_prompt" gorm:"type:text"`
	Active      bool      `json:"active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TokenBalance model
type TokenBalance struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AIEngine        string    `json:"ai_engine" gorm:"uniqueIndex;type:varchar(100)"`
	TotalTokens     int       `json:"total_tokens"`
	UsedTokens      int       `json:"used_tokens" gorm:"default:0"`
	RemainingTokens int       `json:"remaining_tokens"`
	LastReset       time.Time `json:"last_reset"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ConnectedPlatform model
type ConnectedPlatform struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Platform      string    `json:"platform" gorm:"uniqueIndex;type:varchar(100)"`
	Token         string    `json:"token" gorm:"type:text"`
	ClientID      string    `json:"client_id" gorm:"type:varchar(255)"`
	ClientSecret  string    `json:"client_secret" gorm:"type:text"`
	WebhookURL    string    `json:"webhook_url" gorm:"type:text"`
	PhoneNumberID string    `json:"phone_number_id" gorm:"type:varchar(255)"`
	Active        bool      `json:"active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
