package services

import (
	"context"
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"strings"
	"time"
)

type ChatService struct {
	chatRepo       *repository.ChatRepository
	contactService *ContactService
	aiService      *AIService
	productService *ProductService
	logger         *utils.Logger
}

func NewChatService(
	chatRepo *repository.ChatRepository,
	contactService *ContactService,
	aiService *AIService,
	productService *ProductService,
	logger *utils.Logger,
) *ChatService {
	return &ChatService{
		chatRepo:       chatRepo,
		contactService: contactService,
		aiService:      aiService,
		productService: productService,
		logger:         logger,
	}
}

// GetAll returns all chat messages
func (s *ChatService) GetAll() ([]models.ChatMessage, error) {
	return s.chatRepo.FindAll()
}

// GetByID returns a chat message by ID
func (s *ChatService) GetByID(id uint) (*models.ChatMessage, error) {
	return s.chatRepo.FindByID(id)
}

// GetByStatus returns chat messages by status
func (s *ChatService) GetByStatus(status string) ([]models.ChatMessage, error) {
	return s.chatRepo.FindByStatus(status)
}

// GetByContactID returns all messages for a specific contact
func (s *ChatService) GetByContactID(contactID uint) ([]models.ChatMessage, error) {
	return s.chatRepo.FindByContactID(contactID)
}

// GetByChannel returns messages filtered by channel
func (s *ChatService) GetByChannel(channel string) ([]models.ChatMessage, error) {
	return s.chatRepo.FindByChannel(channel)
}

// Create creates a new chat message
func (s *ChatService) Create(message *models.ChatMessage) error {
	s.logger.Info("Creating chat message", "contact_id", message.ContactID)
	return s.chatRepo.Create(message)
}

// Update updates a chat message
func (s *ChatService) Update(message *models.ChatMessage) error {
	s.logger.Info("Updating chat message", "id", message.ID)
	return s.chatRepo.Update(message)
}

// Assign assigns a chat to an agent
func (s *ChatService) Assign(id uint, assignedTo, assignedAgent string) error {
	message, err := s.chatRepo.FindByID(id)
	if err != nil {
		return err
	}

	message.AssignedTo = assignedTo
	message.AssignedAgent = assignedAgent
	message.Status = "Assigned"

	s.logger.Info("Assigning chat", "id", id, "to", assignedTo)
	return s.chatRepo.Update(message)
}

// Resolve marks a chat as resolved
func (s *ChatService) Resolve(id uint) error {
	message, err := s.chatRepo.FindByID(id)
	if err != nil {
		return err
	}

	message.Status = "Resolved"

	s.logger.Info("Resolving chat", "id", id)
	return s.chatRepo.Update(message)
}

// TakeOver - Human agent takes over from AI
func (s *ChatService) TakeOver(id uint, agentName string) error {
	message, err := s.chatRepo.FindByID(id)
	if err != nil {
		return err
	}

	message.AssignedTo = "Human"
	message.AssignedAgent = agentName
	message.Status = "Assigned"

	s.logger.Info("Human agent taking over chat", "id", id, "agent", agentName)
	return s.chatRepo.Update(message)
}

// BackToAI - Return chat back to AI
func (s *ChatService) BackToAI(id uint) error {
	message, err := s.chatRepo.FindByID(id)
	if err != nil {
		return err
	}

	message.AssignedTo = "AI Bot"
	message.AssignedAgent = "AI Assistant"
	message.Status = "Unassigned"

	s.logger.Info("Returning chat to AI", "id", id)
	return s.chatRepo.Update(message)
}

// AddLabel adds a label to a chat message
func (s *ChatService) AddLabel(id uint, labelID string) error {
	message, err := s.chatRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Add label to comma-separated list
	if message.Labels == "" {
		message.Labels = labelID
	} else {
		// Check if label already exists
		labels := strings.Split(message.Labels, ",")
		for _, l := range labels {
			if strings.TrimSpace(l) == labelID {
				return nil // Label already exists
			}
		}
		message.Labels = message.Labels + "," + labelID
	}

	s.logger.Info("Adding label to chat", "id", id, "label", labelID)
	return s.chatRepo.Update(message)
}

// GetStats returns chat statistics
func (s *ChatService) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count by status
	unassigned, _ := s.chatRepo.CountByStatus("Unassigned")
	assigned, _ := s.chatRepo.CountByStatus("Assigned")
	resolved, _ := s.chatRepo.CountByStatus("Resolved")

	stats["unassigned"] = unassigned
	stats["assigned"] = assigned
	stats["resolved"] = resolved
	stats["total"] = unassigned + assigned + resolved

	// Count by channel
	whatsapp, _ := s.chatRepo.CountByChannel("WhatsApp")
	instagram, _ := s.chatRepo.CountByChannel("Instagram")
	telegram, _ := s.chatRepo.CountByChannel("Telegram")

	stats["whatsapp"] = whatsapp
	stats["instagram"] = instagram
	stats["telegram"] = telegram

	// Total tokens used
	totalTokens, _ := s.chatRepo.GetTotalTokens()
	stats["total_tokens"] = totalTokens

	return stats, nil
}

// ProcessIncomingMessage processes an incoming message and generates AI response with RAG
func (s *ChatService) ProcessIncomingMessage(
	platform, phone, name, message string,
) (*models.ChatMessage, string, error) {

	s.logger.Info("💬 Processing incoming message",
		"platform", platform,
		"phone", phone,
		"name", name,
	)

	// 1. Get or create contact
	contact, err := s.contactService.GetOrCreateByChannelID(platform, phone, name)
	if err != nil {
		s.logger.Error("Failed to get/create contact", "error", err)
		return nil, "", err
	}

	s.logger.Info("✅ Contact ready",
		"contact_id", contact.ID,
		"code", contact.Code,
		"name", contact.Name,
	)

	// 2. Save incoming message
	incomingMsg := &models.ChatMessage{
		ContactID:   contact.ID,
		ContactName: contact.Name,
		Message:     message,
		Channel:     platform,
		Status:      "Unassigned",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.chatRepo.Create(incomingMsg); err != nil {
		s.logger.Error("Failed to save incoming message", "error", err)
		return nil, "", err
	}
	s.logger.Info("📝 Incoming message saved", "msg_id", incomingMsg.ID)

	// 3. Generate AI response with RAG context
	ctx := context.Background()
	s.logger.Info("🤖 Generating AI response with RAG...")

	aiResponse, err := s.aiService.GenerateResponse(
		ctx,
		message,
		contact.Name,
		contact.ID, // ✅ Pass contactID for chat history
	)

	if err != nil {
		s.logger.Error("AI processing failed", "error", err)
		aiResponse = "Maaf, saat ini sistem sedang sibuk. Tim kami akan segera menghubungi Anda. 🙏"
	} else {
		s.logger.Info("✅ AI response generated successfully")
	}

	// 4. Save outgoing message
	outgoingMsg := &models.ChatMessage{
		ContactID:     contact.ID,
		ContactName:   contact.Name,
		Message:       message,
		Response:      aiResponse,
		Channel:       platform,
		Status:        "Answered",
		TokensUsed:    0, // Can be calculated if needed
		AssignedTo:    "AI Bot",
		AssignedAgent: "AI Assistant",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.chatRepo.Create(outgoingMsg); err != nil {
		s.logger.Error("Failed to save outgoing message", "error", err)
		return nil, aiResponse, err
	}
	s.logger.Info("✅ Outgoing message saved", "msg_id", outgoingMsg.ID)

	s.logger.Info("🎉 Message processing completed successfully")
	return outgoingMsg, aiResponse, nil
}
