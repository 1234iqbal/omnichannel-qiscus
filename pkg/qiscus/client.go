package qiscus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	appID      string
	secretKey  string
	httpClient *http.Client
}

type Config struct {
	BaseURL   string
	AppID     string
	SecretKey string
	Timeout   time.Duration
}

type AssignAgentRequest struct {
	RoomID  string `json:"room_id"`
	AgentID string `json:"agent_id"`
}

type AssignAgentResponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type GetAgentsResponse struct {
	Status int `json:"status"`
	Data   struct {
		Agents []QiscusAgent `json:"agents"`
	} `json:"data"`
}

type QiscusAgent struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	IsAvailable bool   `json:"is_available"`
}

func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.BaseURL == "" {
		config.BaseURL = "https://omnichannel.qiscus.com"
	}

	return &Client{
		baseURL:   config.BaseURL,
		appID:     config.AppID,
		secretKey: config.SecretKey,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (c *Client) AssignAgent(ctx context.Context, roomID, agentID string) (*AssignAgentResponse, error) {
	endpoint := "/api/v2/admin/service/assign_agent"

	// Prepare form data
	data := url.Values{}
	data.Set("room_id", roomID)
	data.Set("agent_id", agentID)

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Qiscus-App-Id", "bgwus-6yz8hekllunun4v")
	req.Header.Set("Qiscus-Secret-Key", "1d2dfa68b6aef75b8407fc89a532b143")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var assignResponse AssignAgentResponse
	if err := json.Unmarshal(body, &assignResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &assignResponse, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, assignResponse.Message)
	}

	return &assignResponse, nil
}

func (c *Client) GetAgents(ctx context.Context) (*GetAgentsResponse, error) {
	endpoint := "/api/v2/admin/agents"

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Qiscus-App-Id", "bgwus-6yz8hekllunun4v")
	req.Header.Set("Qiscus-Secret-Key", "1d2dfa68b6aef75b8407fc89a532b143")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var agentsResponse GetAgentsResponse
	if err := json.Unmarshal(body, &agentsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &agentsResponse, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return &agentsResponse, nil
}
