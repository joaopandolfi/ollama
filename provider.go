package ollama

import (
	"context"
	"fmt"

	"github.com/agent-api/core"
	"github.com/agent-api/core/message"
	"github.com/agent-api/core/tool"
	"github.com/agent-api/ollama-provider/client"
	"github.com/go-logr/logr"
)

// OllamaProvider implements the LLMProvider interface for Ollama
type OllamaProvider struct {
	host  string
	port  int
	model string

	// client is the internal Ollama HTTP client
	client *client.OllamaClient

	logger logr.Logger
}

// NewOllamaProvider creates a new Ollama provider
//
// TODO:
// - need to handle base URL better (with trailing slashes, etc.)
// - need to construct actual URL using baseURL, port, etc.
func NewOllamaProvider(logger logr.Logger, baseURL string, port int, model string) *OllamaProvider {
	client := client.NewClient(
		client.WithBaseURL("http://localhost:11434/api"),
	)

	return &OllamaProvider{
		host:   baseURL,
		port:   port,
		model:  model,
		client: client,
		logger: logger,
	}
}

func (p *OllamaProvider) GetCapabilities(ctx context.Context) (*core.ProviderCapabilities, error) {
	println("NOT IMPLEMENTED YET")
	return nil, nil
}

// GenerateResponse implements the LLMProvider interface for basic responses
func (p *OllamaProvider) GenerateResponse(ctx context.Context, messages []message.Message) (*message.Message, error) {
	ollamaMessages := convertManyMessagesToOllamaMessages(messages)

	resp, err := p.client.Chat(ctx, &client.ChatRequest{
		Model:    "llama3.2",
		Messages: ollamaMessages,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling client chat method: %w", err)
	}

	return &message.Message{
		Role:    "ai",
		Content: resp.Message.Content,
	}, nil
}

// GenerateWithTools implements the LLMProvider interface for tool-using responses
//
// TODO:
// - handle automatically generating a system message descriptor
func (p *OllamaProvider) GenerateWithTools(ctx context.Context, messages []message.Message, tools []tool.Tool) (*message.Message, error) {
	// Convert tools into Ollama's format
	ollamaTools := make([]client.Tool, len(tools))

	for i, t := range tools {
		ollamaTools[i] = client.Tool{
			Type: "function",
			Function: client.ToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.JSONSchema,
			},
		}
	}

	fmt.Printf("Available tool: \n%s\n%s\n----------\n", ollamaTools[0].Function.Parameters, ollamaTools[0].Function.Name)

	// Convert messages to Ollama format
	ollamaMessages := convertManyMessagesToOllamaMessages(messages)

	// Make the chat request
	resp, err := p.client.Chat(ctx, &client.ChatRequest{
		Model:    p.model,
		Messages: ollamaMessages,
		Tools:    ollamaTools,
	})
	if err != nil {
		return nil, fmt.Errorf("chat request failed: %w", err)
	}

	// Check if the response contains tool calls
	if len(resp.Message.ToolCalls) > 0 {
		// Handle each tool call
		for _, toolCall := range resp.Message.ToolCalls {
			fmt.Printf("LLM tool call: \n%s\n%s\n----------\n", toolCall.Function.Name, toolCall.Function.Arguments)

			// Find the corresponding tool
			var toolToCall *tool.Tool
			for _, t := range tools {
				if t.Name == toolCall.Function.Name {
					toolToCall = &t
					break
				}
			}

			if toolToCall == nil {
				return nil, fmt.Errorf("tool %s not found", toolCall.Function.Name)
			}

			// Call the tool
			result, err := toolToCall.Function(ctx, []byte(toolCall.Function.Arguments))
			if err != nil {
				return nil, fmt.Errorf("tool execution failed: %w", err)
			}

			// Add the tool response to messages
			toolResponseMsg := client.Message{
				Role:    "tool",
				Content: fmt.Sprintf("%v", result),
			}
			ollamaMessages = append(ollamaMessages, &toolResponseMsg)
		}
	}

	// Make another chat request with the tool response
	resp, err = p.client.Chat(ctx, &client.ChatRequest{
		Model:    p.model,
		Messages: ollamaMessages,
		Tools:    ollamaTools,
	})
	if err != nil {
		return nil, fmt.Errorf("follow-up chat request failed: %w", err)
	}

	// Return the final response
	return &message.Message{
		Role:    "assistant",
		Content: resp.Message.Content,
	}, nil
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

//func (p *OllamaProvider) apiChat(ctx context.Context, jsonBody []byte) (*ollamaResponse, error) {
//req, err := http.NewRequestWithContext(ctx, "POST", p.host+":"+strconv.Itoa(p.port)+"/api/chat", bytes.NewBuffer(jsonBody))
//if err != nil {
//return nil, fmt.Errorf("error creating request: %w", err)
//}
//req.Header.Set("Content-Type", "application/json")

//resp, err := p.client.Do(req)
//if err != nil {
//return nil, fmt.Errorf("error making request: %w", err)
//}
//defer resp.Body.Close()

//body, err := io.ReadAll(resp.Body)
//if err != nil {
//return nil, fmt.Errorf("error reading response: %w", err)
//}

//if resp.StatusCode != http.StatusOK {
//return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
//}

//var ollamaResp ollamaResponse
//if err := json.Unmarshal(body, &ollamaResp); err != nil {
//return nil, fmt.Errorf("error unmarshaling response: %w", err)
//}

//if ollamaResp.Error != "" {
//return nil, fmt.Errorf("ollama error: %s", ollamaResp.Error)
//}

//return &ollamaResp, nil
//}
