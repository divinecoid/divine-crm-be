package models

import (
	"time"
)

// ==================== CONTACTS ====================

type Contact struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Code          string    `json:"code" gorm:"unique;not null"`
	Channel       string    `json:"channel"`    // WhatsApp, Telegram, Instagram
	ChannelID     string    `json:"channel_id"` // Phone number, username, etc
	Name          string    `json:"name"`
	ContactStatus string    `json:"contact_status"` // Leads, Contact
	Temperature   string    `json:"temperature"`    // Cold, Warm, Hot
	FirstContact  time.Time `json:"first_contact"`
	LastContact   time.Time `json:"last_contact"`
	LastAgent     string    `json:"last_agent"`
	LastAgentType string    `json:"last_agent_type"` // AI, Human
	Notes         string    `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ==================== PRODUCTS ====================

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Code        string    `json:"code" gorm:"unique;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Description string    `json:"description" gorm:"type:text"`
	UploadedBy  string    `json:"uploaded_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ==================== CHAT LABELS ====================

type ChatLabel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Label       string    `json:"label" gorm:"unique;not null"`
	Description string    `json:"description"`
	Color       string    `json:"color"` // Red, Purple, Pink, Green, etc
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ==================== CHAT MESSAGES ====================

type ChatMessage struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ContactID     uint      `json:"contact_id"`
	Contact       Contact   `json:"contact" gorm:"foreignKey:ContactID"`
	ContactName   string    `json:"contact_name"`
	Message       string    `json:"message" gorm:"type:text"`
	Response      string    `json:"response" gorm:"type:text"`
	Status        string    `json:"status"` // Unassigned, Pending, Assigned, Resolved
	AssignedTo    string    `json:"assigned_to"`
	AssignedAgent string    `json:"assigned_agent"`
	Channel       string    `json:"channel"`
	Labels        string    `json:"labels"` // Comma-separated label IDs
	TokensUsed    int       `json:"tokens_used"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ==================== AI CONFIGURATION ====================

type AIConfiguration struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AIEngine  string    `json:"ai_engine" gorm:"unique;not null"` // openai, deepseek, grok, gemini
	Token     string    `json:"token"`
	Endpoint  string    `json:"endpoint"`
	Model     string    `json:"model"`
	Active    bool      `json:"active" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ==================== CONNECTED PLATFORMS ====================

type ConnectedPlatform struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Platform      string    `json:"platform" gorm:"unique;not null"` // WhatsApp, Instagram, Telegram
	PlatformID    string    `json:"platform_id"`                     // Phone, Username, Bot ID
	Token         string    `json:"token"`
	PhoneNumberID string    `json:"phone_number_id"` // For WhatsApp
	ClientID      string    `json:"client_id"`
	ClientSecret  string    `json:"client_secret"`
	WebhookURL    string    `json:"webhook_url"`
	Active        bool      `json:"active" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ==================== AI AGENTS ====================

type AIAgent struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"` // Diva, Clara, Kana, Gema
	AIEngine    string    `json:"ai_engine"`                   // openai, deepseek, grok, gemini
	BasicPrompt string    `json:"basic_prompt" gorm:"type:text"`
	Active      bool      `json:"active" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ==================== HUMAN AGENTS ====================

type HumanAgent struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Username    string     `json:"username" gorm:"unique;not null"`
	Password    string     `json:"-"` // Hidden from JSON
	Email       string     `json:"email"`
	FullName    string     `json:"full_name"`
	Role        string     `json:"role"` // Admin, Agent, Supervisor
	Active      bool       `json:"active" gorm:"default:true"`
	LatestLogin *time.Time `json:"latest_login"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ==================== BROADCAST TEMPLATE ====================

type BroadcastTemplate struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text"`
	Channel   string    `json:"channel"` // WhatsApp, Instagram, Telegram, All
	CreatedBy string    `json:"created_by"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ==================== BROADCAST HISTORY ====================

type BroadcastHistory struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	TemplateID  uint              `json:"template_id"`
	Template    BroadcastTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	SentTo      int               `json:"sent_to"` // Total recipients
	Successful  int               `json:"successful"`
	Failed      int               `json:"failed"`
	Status      string            `json:"status"` // Pending, Processing, Completed, Failed
	SentBy      string            `json:"sent_by"`
	CreatedAt   time.Time         `json:"created_at"`
	CompletedAt *time.Time        `json:"completed_at"`
}

// ==================== QUICK REPLY ====================

type QuickReply struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Trigger   string    `json:"trigger" gorm:"unique;not null"` // Keyword to trigger
	Response  string    `json:"response" gorm:"type:text"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ==================== API SETTINGS ====================

type APISettings struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Key         string    `json:"key" gorm:"unique;not null"`
	Value       string    `json:"value" gorm:"type:text"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ==================== TOKEN BALANCE ====================

type TokenBalance struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"user_id"`
	AIEngine        string    `json:"ai_engine"`
	TotalTokens     int64     `json:"total_tokens"`
	UsedTokens      int64     `json:"used_tokens"`
	RemainingTokens int64     `json:"remaining_tokens"`
	LastReset       time.Time `json:"last_reset"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ==================== ANALYTICS ====================

type Analytics struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Date              time.Time `json:"date" gorm:"type:date"`
	TotalMessages     int       `json:"total_messages"`
	UnassignedChats   int       `json:"unassigned_chats"`
	AssignedChats     int       `json:"assigned_chats"`
	ResolvedChats     int       `json:"resolved_chats"`
	NewContacts       int       `json:"new_contacts"`
	TotalTokensUsed   int       `json:"total_tokens_used"`
	WhatsAppMessages  int       `json:"whatsapp_messages"`
	InstagramMessages int       `json:"instagram_messages"`
	TelegramMessages  int       `json:"telegram_messages"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
