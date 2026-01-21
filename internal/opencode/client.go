package opencode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is an HTTP client for the OpenCode server API
type Client struct {
	serverURL  string
	username   string
	password   string
	modelID    string
	httpClient *http.Client
}

// Config holds configuration for the OpenCode client
type Config struct {
	ServerURL string
	Username  string
	Password  string
	ModelID   string
	Timeout   time.Duration
}

// NewClient creates a new OpenCode API client
func NewClient(cfg Config) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	return &Client{
		serverURL: strings.TrimRight(cfg.ServerURL, "/"),
		username:  cfg.Username,
		password:  cfg.Password,
		modelID:   cfg.ModelID,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// CreateSessionRequest is the request body for creating a new session
type CreateSessionRequest struct {
	ModelID string `json:"model_id,omitempty"`
}

// CreateSessionResponse is the response from creating a new session
type CreateSessionResponse struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

// MessagePart represents a part of a message
type MessagePart struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// SendMessageRequest is the request body for sending a message
type SendMessageRequest struct {
	Parts   []MessagePart `json:"parts"`
	ModelID string        `json:"model_id,omitempty"`
}

// ResponsePart represents a part of the response message
type ResponsePart struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Text      string `json:"text,omitempty"`
	SessionID string `json:"sessionID"`
	MessageID string `json:"messageID"`
}

// MessageInfo contains metadata about the response message
type MessageInfo struct {
	ID        string `json:"id"`
	SessionID string `json:"sessionID"`
	Role      string `json:"role"`
	ModelID   string `json:"modelID"`
}

// SendMessageResponse is the response from sending a message
type SendMessageResponse struct {
	Info  MessageInfo    `json:"info"`
	Parts []ResponsePart `json:"parts"`
}

// Content extracts the text content from the response parts
func (r *SendMessageResponse) Content() string {
	for _, part := range r.Parts {
		if part.Type == "text" && part.Text != "" {
			return part.Text
		}
	}
	return ""
}

// SessionID returns the session ID from the response
func (r *SendMessageResponse) SessionID() string {
	return r.Info.SessionID
}

// CreateSession creates a new OpenCode session
func (c *Client) CreateSession() (string, error) {
	url := c.serverURL + "/session"

	reqBody := CreateSessionRequest{
		ModelID: c.modelID,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create session: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result CreateSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.ID, nil
}

// SendMessage sends a message to an existing session
func (c *Client) SendMessage(sessionID, content string) (*SendMessageResponse, error) {
	url := fmt.Sprintf("%s/session/%s/message", c.serverURL, sessionID)

	reqBody := SendMessageRequest{
		Parts: []MessagePart{
			{Type: "text", Text: content},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to send message: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// setHeaders sets common headers including basic auth
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}
}

// ServerURL returns the configured server URL
func (c *Client) ServerURL() string {
	return c.serverURL
}

// ModelID returns the configured model ID
func (c *Client) ModelID() string {
	return c.modelID
}
