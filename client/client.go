package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OllamaClient represents an Ollama API client
type OllamaClient struct {
	baseURL    string
	httpClient *http.Client
}

// ClientOption is a function that modifies the client
type ClientOption func(*OllamaClient)

// NewClient creates a new Ollama client
func NewClient(options ...ClientOption) *OllamaClient {
	client := &OllamaClient{
		baseURL:    "http://localhost:11434/api",
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// WithBaseURL sets the base URL for the client
func WithBaseURL(url string) ClientOption {
	return func(c *OllamaClient) {
		c.baseURL = strings.TrimRight(url, "/")
	}
}

// WithHTTPClient sets the HTTP client for the API client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *OllamaClient) {
		c.httpClient = httpClient
	}
}

// Chat sends a chat request to the Ollama API
func (c *OllamaClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	return c.ChatStream(ctx, req, nil)
}

// ChatStream sends a chat request and handles streaming responses
func (c *OllamaClient) ChatStream(ctx context.Context, req *ChatRequest, handler func(*ChatResponse) error) (*ChatResponse, error) {
	if req.Stream == nil {
		stream := false
		req.Stream = &stream
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	decoder := json.NewDecoder(resp.Body)
	var lastResponse *ChatResponse

	for {
		var response ChatResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		lastResponse = &response

		if handler != nil {
			if err := handler(&response); err != nil {
				return nil, fmt.Errorf("stream handler failed: %w", err)
			}
		}

		if response.Done {
			break
		}
	}

	return lastResponse, nil
}
