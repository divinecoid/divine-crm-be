package services

import (
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"divine-crm/internal/utils"
	"time"
)

// ChatService handles chat business logic
type ChatService struct {
	chatRepo       *repository.ChatRepository
	contactService *ContactService
	aiService      *AIService
	logger         *utils.Logger
}

// NewChatService creates a new chat service
func NewChatService(
	chatRepo *repository.ChatRepository,
	contactService *ContactService,
	aiService *AIService,
	logger *utils.Logger,
) *ChatService {
	return &ChatService{
		chatRepo:       chatRepo,
		contactService: contactService,
		aiService:      aiService,
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

// ProcessIncomingMessage processes an incoming message and generates AI response
func (s *ChatService) ProcessIncomingMessage(
	platform, phone, name, message string,
) (*models.ChatMessage, string, error) {

	s.logger.Info("ChatService: Starting to process incoming message")

	// 1. Get or create contact using ContactService method
	s.logger.Info("ChatService: Getting or creating contact...", "phone", phone, "name", name)

	// ✅ Use the GetOrCreateByChannelID method from ContactService
	contact, err := s.contactService.GetOrCreateByChannelID(platform, phone, name)
	if err != nil {
		s.logger.Error("ChatService: Failed to get/create contact", "error", err)
		return nil, "", err
	}

	s.logger.Info("ChatService: Contact ready", "contact_id", contact.ID, "code", contact.Code)

	// 2. Save incoming message
	s.logger.Info("ChatService: Saving incoming message...")
	incomingMsg := &models.ChatMessage{
		ContactID:   contact.ID,
		ContactName: name,
		Message:     message,
		Channel:     platform,
		Status:      "Unassigned",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.chatRepo.Create(incomingMsg); err != nil {
		s.logger.Error("ChatService: Failed to save incoming message", "error", err)
		return nil, "", err
	}
	s.logger.Info("ChatService: Incoming message saved", "msg_id", incomingMsg.ID)

	// 3. Get AI response
	s.logger.Info("ChatService: Requesting AI response...")

	// ✅ Use ProcessWithActiveAgent (no agentID parameter needed)
	aiResponse, tokens, err := s.aiService.ProcessWithActiveAgent(message)
	if err != nil {
		s.logger.Error("ChatService: AI processing failed", "error", err)
		// Return a default response instead of failing
		aiResponse = "Maaf, saat ini sistem sedang sibuk. Tim kami akan segera menghubungi Anda."
		tokens = 0
	}
	s.logger.Info("ChatService: AI response received", "tokens", tokens)

	// 4. Save outgoing message (with AI response)
	s.logger.Info("ChatService: Saving outgoing message...")
	outgoingMsg := &models.ChatMessage{
		ContactID:     contact.ID,
		ContactName:   name,
		Message:       message,    // Original user message
		Response:      aiResponse, // AI response
		Channel:       platform,
		Status:        "Answered",
		TokensUsed:    tokens,
		AssignedTo:    "AI Bot",
		AssignedAgent: "AI Assistant",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.chatRepo.Create(outgoingMsg); err != nil {
		s.logger.Error("ChatService: Failed to save outgoing message", "error", err)
		// Still return AI response even if save fails
		return nil, aiResponse, err
	}
	s.logger.Info("ChatService: Outgoing message saved", "msg_id", outgoingMsg.ID)

	s.logger.Info("ChatService: Message processing completed successfully")
	return outgoingMsg, aiResponse, nil
}
