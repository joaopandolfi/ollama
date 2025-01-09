package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/agent-api/core"
	"github.com/agent-api/core/message"
	"github.com/agent-api/core/tool"
	"github.com/go-logr/logr"
)

// OllamaProvider implements the LLMProvider interface for Ollama
type OllamaProvider struct {
	host   string
	port   int
	model  string
	client *http.Client

	logger logr.Logger
}

type ollamaRequest struct {
	Model    string           `json:"model"`
	Messages []*ollamaMessage `json:"messages"`
	Stream   bool             `json:"stream"`
}

type ollamaResponse struct {
	Message ollamaMessage `json:"message"`
	Error   string        `json:"error,omitempty"`
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(logger logr.Logger, baseURL string, port int, model string) *OllamaProvider {
	// todo - fix this
	//if !strings.HasSuffix(baseURL, "/") {
	//baseURL += "/"
	//}

	return &OllamaProvider{
		host:   baseURL,
		port:   port,
		model:  model,
		client: &http.Client{},
		logger: logger,
	}
}

func GetCapabilities(ctx context.Context) (*core.ProviderCapabilities, error) {
	println("NOT IMPLEMENTED YET")
	return nil, nil
}

// GenerateResponse implements the LLMProvider interface for basic responses
func (p *OllamaProvider) GenerateResponse(ctx context.Context, messages []message.Message) (*message.Message, error) {
	ollamaMessages := convertManyMessagesToOllamaMessages(messages)

	reqBody := ollamaRequest{
		Model:    p.model,
		Messages: ollamaMessages,
		Stream:   false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	response, err := p.apiChat(ctx, jsonBody)
	if err != nil {
		return nil, err
	}

	return &message.Message{
		Role:    "ai",
		Content: response.Message.Content,
	}, nil
}

// GenerateWithTools implements the LLMProvider interface for tool-using responses
func (p *OllamaProvider) GenerateWithTools(ctx context.Context, messages []message.Message, tools []tool.Tool) (*message.Message, error) {
	// For Ollama, we'll embed tool information in the prompt since it doesn't natively support tools
	// Create a description of available tools
	var toolsDesc strings.Builder
	toolsDesc.WriteString("\nAvailable tools:\n")
	for _, tool := range tools {
		toolsDesc.WriteString(fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description))
	}
	toolsDesc.WriteString("\nTo use a tool, respond with: <tool>tool_name|{\"param\":\"value\"}</tool>")

	// Add tools description to the last user message
	lastMsgIdx := len(messages) - 1
	if lastMsgIdx >= 0 && messages[lastMsgIdx].Role == "user" {
		messages[lastMsgIdx].Content += toolsDesc.String()
	}

	// Use regular generation
	return p.GenerateResponse(ctx, messages)
}

// GenerateStream streams the response token by token
func (p *OllamaProvider) GenerateStream(ctx context.Context, messages []message.Message, opts *core.InferenceOptions) (<-chan *message.Message, <-chan error) {
	println("NOT IMPLEMENTED YET")
	return nil, nil
}

// GenerateStreamWithTools streams a response with tools token by token
func (p *OllamaProvider) GenerateStreamWithTools(ctx context.Context, messages []message.Message, tools []tool.Tool, opts *core.InferenceOptions) (<-chan *message.Message, <-chan error) {
	println("NOT IMPLEMENTED YET")
	return nil, nil
}

// ValidatePrompt checks if a prompt is valid for the provider
func (p *OllamaProvider) ValidatePrompt(ctx context.Context, messages []message.Message) error {
	println("NOT IMPLEMENTED YET")
	return nil
}

// EstimateTokens estimates the number of tokens in a message
func (p *OllamaProvider) EstimateTokens(ctx context.Context, message string) (int, error) {
	println("NOT IMPLEMENTED YET")
	return 0, nil
}

// GetModelList returns available models for this provider
func (p *OllamaProvider) GetModelList(ctx context.Context) ([]string, error) {
	println("NOT IMPLEMENTED YET")
	return nil, nil
}

func (p *OllamaProvider) apiChat(ctx context.Context, jsonBody []byte) (*ollamaResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", p.host+":"+strconv.Itoa(p.port)+"/api/chat", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if ollamaResp.Error != "" {
		return nil, fmt.Errorf("ollama error: %s", ollamaResp.Error)
	}

	return &ollamaResp, nil
}
