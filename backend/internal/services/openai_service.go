package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fact-check/internal/config"

	"github.com/sirupsen/logrus"
)

type OpenAIService struct {
	config *config.Config
	logger *logrus.Logger
}

type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

func NewOpenAIService(cfg *config.Config, logger *logrus.Logger) *OpenAIService {
	return &OpenAIService{
		config: cfg,
		logger: logger,
	}
}

func (s *OpenAIService) VerifyNews(content string, link string, photoURL string) (string, string, error) {
	// Check if OpenAI API key is configured
	if s.config.OpenAIAPIKey == "" || s.config.OpenAIAPIKey == "your-openai-api-key" {
		return "uncertain", "OpenAI API not configured. Please configure your OpenAI API key to enable fact-checking.", nil
	}
	
	prompt := s.buildPrompt(content, link, photoURL)

	request := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a fact-checking expert. Analyze the provided news content and determine if it's likely to be true or false. Provide a clear explanation for your assessment.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 500,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal OpenAI request: %w", err)
	}

	req, err := http.NewRequest("POST", s.config.OpenAIEndpoint+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.OpenAIAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("OpenAI API error: %s", string(body))
		
		// Try to parse the error response for better error messages
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		
		if json.Unmarshal(body, &errorResp) == nil {
			switch errorResp.Error.Code {
			case "insufficient_quota":
				return "", "", fmt.Errorf("OpenAI API quota exceeded: %s. Please check your billing and upgrade your plan.", errorResp.Error.Message)
			case "rate_limit_exceeded":
				return "", "", fmt.Errorf("OpenAI API rate limit exceeded: %s. Please try again later.", errorResp.Error.Message)
			default:
				return "", "", fmt.Errorf("OpenAI API error (%s): %s", errorResp.Error.Code, errorResp.Error.Message)
			}
		}
		
		return "", "", fmt.Errorf("OpenAI API returned status %d", resp.StatusCode)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal OpenAI response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", "", fmt.Errorf("no response from OpenAI API")
	}

	response := openAIResp.Choices[0].Message.Content

	// Parse the response to extract status and explanation
	status, explanation := s.parseOpenAIResponse(response)

	return status, explanation, nil
}

func (s *OpenAIService) buildPrompt(content string, link string, photoURL string) string {
	prompt := fmt.Sprintf("Please fact-check the following news content:\n\nContent: %s\n", content)

	if link != "" {
		prompt += fmt.Sprintf("Source Link: %s\n", link)
	}

	if photoURL != "" {
		prompt += fmt.Sprintf("Photo URL: %s\n", photoURL)
	}

	prompt += "\nPlease respond with:\n1. A clear assessment: 'TRUE', 'FALSE', or 'UNCERTAIN'\n2. A detailed explanation for your assessment\n3. Any relevant context or sources you considered"

	return prompt
}

func (s *OpenAIService) parseOpenAIResponse(response string) (string, string) {
	// Simple parsing logic - in production, you might want more sophisticated parsing
	response = response + " " // Add space to ensure we can find the end

	// Look for status indicators
	var status string
	if contains(response, "TRUE") || contains(response, "true") {
		status = "true"
	} else if contains(response, "FALSE") || contains(response, "false") {
		status = "false"
	} else {
		status = "uncertain"
	}

	// The explanation is the full response
	explanation := response

	return status, explanation
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			contains(s[1:], substr))))
}

// IsAvailable checks if the OpenAI service is properly configured and available
func (s *OpenAIService) IsAvailable() bool {
	return s.config.OpenAIAPIKey != "" && s.config.OpenAIAPIKey != "your-openai-api-key"
}

// GetServiceStatus returns the current status of the OpenAI service
func (s *OpenAIService) GetServiceStatus() map[string]interface{} {
	status := map[string]interface{}{
		"available": s.IsAvailable(),
	}
	
	if s.IsAvailable() {
		status["configured"] = true
		status["message"] = "OpenAI service is configured and available"
	} else {
		status["configured"] = false
		status["message"] = "OpenAI API key not configured. Please set OPENAI_API_KEY environment variable."
	}
	
	return status
}
