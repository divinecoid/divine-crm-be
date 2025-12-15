package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Contact struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Code          string    `json:"code" gorm:"uniqueIndex"`
	Channel       string    `json:"channel"`
	ChannelID     string    `json:"channel_id"`
	Name          string    `json:"name"`
	Temperature   string    `json:"temperature"`
	FirstContact  time.Time `json:"first_contact"`
	LastContact   time.Time `json:"last_contact"`
	LastAgent     string    `json:"last_agent"`
	LastAgentType string    `json:"last_agent_type"`
	Notes         string    `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Lead struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Code          string    `json:"code" gorm:"uniqueIndex"`
	Channel       string    `json:"channel"`
	ChannelID     string    `json:"channel_id"`
	Temperature   string    `json:"temperature"`
	FirstContact  time.Time `json:"first_contact"`
	LastContact   time.Time `json:"last_contact"`
	LastAgent     string    `json:"last_agent"`
	LastAgentType string    `json:"last_agent_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Code        string    `json:"code" gorm:"uniqueIndex"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Description string    `json:"description" gorm:"type:text"`
	UploadedBy  string    `json:"uploaded_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ChatMessage struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ContactID     uint      `json:"contact_id"`
	ContactName   string    `json:"contact_name"`
	Message       string    `json:"message" gorm:"type:text"`
	Response      string    `json:"response" gorm:"type:text"`
	Status        string    `json:"status" gorm:"default:'Unassigned'"`
	AssignedTo    string    `json:"assigned_to"`
	AssignedAgent string    `json:"assigned_agent"`
	Channel       string    `json:"channel"`
	Labels        string    `json:"labels" gorm:"type:text"`
	TokensUsed    int       `json:"tokens_used" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ChatLabel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QuickReply struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Trigger   string    `json:"trigger" gorm:"uniqueIndex"`
	Response  string    `json:"response" gorm:"type:text"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BroadcastTemplate struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Content   string    `json:"content" gorm:"type:text"`
	Type      string    `json:"type"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Broadcast struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	TemplateID     uint      `json:"template_id"`
	TemplateName   string    `json:"template_name"`
	Message        string    `json:"message" gorm:"type:text"`
	Channel        string    `json:"channel"`
	TargetFilter   string    `json:"target_filter"`
	TotalSent      int       `json:"total_sent"`
	TotalDelivered int       `json:"total_delivered"`
	TotalFailed    int       `json:"total_failed"`
	Status         string    `json:"status" gorm:"default:pending"`
	SentBy         string    `json:"sent_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AIConfiguration struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AIEngine  string    `json:"ai_engine" gorm:"uniqueIndex"`
	Token     string    `json:"token"`
	Endpoint  string    `json:"endpoint"`
	Model     string    `json:"model"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AIAgent struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"uniqueIndex"`
	AIEngine     string    `json:"ai_engine"`
	Platform     string    `json:"platform" gorm:"default:'All'"`
	BasicPrompt  string    `json:"basic_prompt" gorm:"type:text"`
	IntroMessage string    `json:"intro_message" gorm:"type:text"`
	Active       bool      `json:"active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TokenBalance struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	AIEngine        string    `json:"ai_engine" gorm:"uniqueIndex"`
	TotalTokens     int       `json:"total_tokens"`
	UsedTokens      int       `json:"used_tokens" gorm:"default:0"`
	RemainingTokens int       `json:"remaining_tokens"`
	LastReset       time.Time `json:"last_reset"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ConnectedPlatform struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Platform      string    `json:"platform" gorm:"uniqueIndex"`
	PlatformID    string    `json:"platform_id"`
	Token         string    `json:"token"`
	ClientID      string    `json:"client_id"`
	ClientSecret  string    `json:"client_secret"`
	WebhookURL    string    `json:"webhook_url"`
	PhoneNumberID string    `json:"phone_number_id"`
	Active        bool      `json:"active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AIRequest struct {
	Model       string      `json:"model"`
	Messages    []AIMessage `json:"messages"`
	Temperature float64     `json:"temperature,omitempty"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

type WhatsAppWebhook struct {
	Object string `json:"object"`
	Entry  []struct {
		Changes []struct {
			Value struct {
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From string `json:"from"`
					Text struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

const (
	TempHot  = "Hot"
	TempWarm = "Warm"
	TempCold = "Cold"
)
