package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Model string

const (
	ModelHaiku  Model = "claude-3-haiku-20240307"
	ModelSonnet Model = "claude-3-sonnet-20240229"
	ModelOpus   Model = "claude-3-opus-20240229"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type Metadata struct {
	UserId string `json:"user_id,omitempty"`
}

type MessageRequest struct {
	Model       Model     `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	System      string    `json:"system,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Metadata    Metadata  `json:"metadata,omitempty"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Usage struct {
	Input  int `json:"input_tokens"`
	Output int `json:"output_tokens"`
}

type MessageResponse struct {
	Id      string    `json:"id"`
	Role    Role      `json:"role"`
	Content []Content `json:"content"`
	Model   Model     `json:"model"`
	Usage   Usage     `json:"usage"`
}

func (c *Client) CreateMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := "https://api.anthropic.com/v1/messages"
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request %s: %w", url, err)
	}

	rq.Header.Set("x-api-key", c.apiKey)
	rq.Header.Set("anthropic-version", "2023-06-01")
	rq.Header.Set("content-type", "application/json")

	resp, err := c.httpClient.Do(rq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %s: %d", url, resp.StatusCode)
	}

	rp := MessageResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&rp); err != nil {
		return nil, fmt.Errorf("failed to decode response %s: %w", url, err)
	}

	return &rp, nil
}
