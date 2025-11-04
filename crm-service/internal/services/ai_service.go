package services

import (
	"bytes"
	"divine-crm/internal/config"
	"divine-crm/internal/models"
	"divine-crm/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AIService handles AI processing logic
type AIService struct {
	repo   *repository.AIRepository
	config *config.Config
}

// NewAIService creates a new AI service
func NewAIService(repo *repository.AIRepository, cfg *config.Config) *AIService {
	return &AIService{
		repo:   repo,
		config: cfg,
	}
}

// AIRequest represents an AI API request
type AIRequest struct {
	Model       string      `json:"model"`
	Messages    []AIMessage `json:"messages"`
	Temperature float64     `json:"temperature,omitempty"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
}

// AIMessage represents a message in AI request
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIResponse represents an AI API response
type AIResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// ProcessMessage processes a message with AI
func (s *AIService) ProcessMessage(message string, agentID uint) (string, int, error) {
	// Get agent
	agent, err := s.repo.GetAgentByID(agentID)
	if err != nil {
		return "", 0, fmt.Errorf("agent not found: %w", err)
	}

	// Get AI configuration
	config, err := s.repo.GetConfiguration(agent.AIEngine)
	if err != nil {
		return "", 0, fmt.Errorf("AI configuration not found: %w", err)
	}

	// Call AI API
	return s.callAI(config, agent.BasicPrompt, message)
}

// ProcessWithActiveAgent processes a message with the active agent
func (s *AIService) ProcessWithActiveAgent(message string) (string, int, error) {
	// Get active agent
	agent, err := s.repo.GetActiveAgent()
	if err != nil {
		return "", 0, fmt.Errorf("no active AI agent found: %w", err)
	}

	// Get AI configuration
	config, err := s.repo.GetConfiguration(agent.AIEngine)
	if err != nil {
		return "", 0, fmt.Errorf("AI configuration not found: %w", err)
	}

	// Call AI API
	return s.callAI(config, agent.BasicPrompt, message)
}

func (s *AIService) callAI(config *models.AIConfiguration, systemPrompt, userMessage string) (string, int, error) {
	// Prepare request
	aiReq := AIRequest{
		Model: config.Model,
		Messages: []AIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		Temperature: 0.7,
		MaxTokens:   500,
	}

	reqBody, err := json.Marshal(aiReq)
	if err != nil {
		return "", 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", config.Endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	// Send request
	client := &http.Client{Timeout: s.config.AI.RequestTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("AI API error [%d]: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var aiResp AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return "", 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(aiResp.Choices) == 0 {
		return "", 0, fmt.Errorf("no response from AI")
	}

	return aiResp.Choices[0].Message.Content, aiResp.Usage.TotalTokens, nil
}
