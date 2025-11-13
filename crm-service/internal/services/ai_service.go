package services

import (
	"context"
	"divine-crm/internal/config"
	"divine-crm/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type AIService struct {
	repo          *repository.AIRepository
	vectorService *VectorService // ✅ Add this
	config        *config.Config
	client        *http.Client
}

func NewAIService(repo *repository.AIRepository, vectorService *VectorService, cfg *config.Config) *AIService {
	return &AIService{
		repo:          repo,
		vectorService: vectorService,
		config:        cfg,
		client:        &http.Client{Timeout: 30 * time.Second},
	}
}

// GenerateResponse generates AI response with RAG
func (s *AIService) GenerateResponse(ctx context.Context, userMessage string, contactName string, contactID uint) (string, error) {
	// 1. Build RAG context from vector search
	ragContext := ""
	if s.vectorService != nil {
		if context, err := s.vectorService.BuildRAGContext(userMessage); err == nil {
			ragContext = context
		}
	}

	// 2. Build system prompt with RAG context
	systemPrompt := s.buildSystemPrompt(ragContext)

	// 3. Generate response with OpenAI
	response, err := s.generateWithOpenAI(ctx, systemPrompt, userMessage, contactName)
	if err != nil {
		return "", err
	}

	// 4. Save chat history with embeddings (async)
	if s.vectorService != nil {
		go func() {
			_ = s.vectorService.SaveChatWithEmbedding(
				contactID,
				userMessage,
				response,
				"neutral", // Can be enhanced with sentiment analysis
				"inquiry", // Can be enhanced with intent classification
			)
		}()
	}

	return response, nil
}

func (s *AIService) buildSystemPrompt(ragContext string) string {
	basePrompt := `Anda adalah AI assistant untuk Divine CRM, sebuah platform CRM berbasis AI.

Tugas Anda:
- Jawab pertanyaan customer dengan ramah, profesional, dan helpful
- Gunakan informasi dari knowledge base yang disediakan
- Jika customer menanyakan harga atau fitur, jelaskan dengan detail
- Jika tidak yakin atau tidak ada info di knowledge base, arahkan ke tim support
- Selalu akhiri dengan tawaran bantuan lebih lanjut
- Gunakan emoji secukupnya untuk membuat respons lebih friendly

Tone: Ramah, profesional, helpful, natural (seperti chat biasa)`

	if ragContext != "" {
		basePrompt += "\n\n" + ragContext
		basePrompt += "\n\nGunakan informasi di atas untuk menjawab pertanyaan customer."
	}

	return basePrompt
}

func (s *AIService) generateWithOpenAI(ctx context.Context, systemPrompt, userMessage, contactName string) (string, error) {
	reqBody := map[string]interface{}{
		"model": s.config.OpenAI.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("%s bertanya: %s", contactName, userMessage),
			},
		},
		"temperature": 0.7,
		"max_tokens":  500,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.OpenAI.Endpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.OpenAI.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	content := message["content"].(string)

	return strings.TrimSpace(content), nil
}
