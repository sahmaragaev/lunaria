package services

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sahmaragaev/lunaria-backend/internal/config"
)

type GrokService struct {
	client *resty.Client
	config *config.GrokConfig
}

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GrokRequest struct {
	Model       string       `json:"model"`
	Messages    []LLMMessage `json:"messages"`
	MaxTokens   int          `json:"max_tokens"`
	Temperature float64      `json:"temperature"`
	Stream      bool         `json:"stream"`
}

type GrokResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func NewGrokService(cfg *config.GrokConfig) *GrokService {
	client := resty.New()
	client.SetHeader("Authorization", "Bearer "+cfg.APIKey)
	client.SetHeader("Content-Type", "application/json")
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.x.ai/v1/chat/completions"
	}

	return &GrokService{
		client: client,
		config: cfg,
	}
}

func (g *GrokService) SendMessage(ctx context.Context, messages []LLMMessage) (string, error) {
	request := GrokRequest{
		Model:       g.config.Model,
		Messages:    messages,
		MaxTokens:   g.config.MaxTokens,
		Temperature: g.config.Temperature,
		Stream:      false,
	}

	var response GrokResponse

	resp, err := g.client.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&response).
		Post(g.config.BaseURL)

	if err != nil {
		return "", fmt.Errorf("failed to send request to Grok: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("Grok API returned status %d: %s", resp.StatusCode(), resp.String())
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from Grok")
	}

	return response.Choices[0].Message.Content, nil
}

func (g *GrokService) SendMiniMessage(ctx context.Context, messages []LLMMessage) (string, error) {
	request := GrokRequest{
		Model:       g.config.MiniModel,
		Messages:    messages,
		MaxTokens:   g.config.MaxTokens,
		Temperature: 0.7,
		Stream:      false,
	}

	var response GrokResponse

	resp, err := g.client.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&response).
		Post(g.config.BaseURL)

	if err != nil {
		return "", fmt.Errorf("failed to send request to Grok Mini: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("Grok Mini API returned status %d: %s", resp.StatusCode(), resp.String())
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from Grok Mini")
	}

	return response.Choices[0].Message.Content, nil
}
