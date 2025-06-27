package qiscus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"qiscus-agent-allocation/internal/domain/entity"
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

func (c *Client) GetAgents() ([]entity.QiscusAgent, error) {
	url := "/api/v2/admin/agents"

	req, err := http.NewRequest("GET", c.baseURL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Qiscus-App-Id", c.appID)
	req.Header.Set("Qiscus-Secret-Key", c.secretKey)

	// Make HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var response entity.GetAgentsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response.Data.Agents, nil
}

func (c *Client) AssignAgent(roomID, agentID string) error {
	url := "/api/v1/admin/service/assign_agent"

	// Prepare request body
	requestBody := entity.AssignAgentRequest{
		RoomID:  roomID,
		AgentID: agentID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Qiscus-App-Id", c.appID)
	req.Header.Set("Qiscus-Secret-Key", c.secretKey)

	// Make HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
