package ollama

import (
	"context"
	"fmt"

	"github.com/agent-api/core"
	"github.com/agent-api/core/message"
	"github.com/agent-api/core/model"
	"github.com/agent-api/ollama-provider/client"
	"github.com/go-logr/logr"
)

// Provider implements the LLMProvider interface for Ollama
type Provider struct {
	host string
	port int

	model *model.Model

	// client is the internal Ollama HTTP client
	client *client.OllamaClient

	logger logr.Logger
}

type ProviderOpts struct {
	BaseURL string
	Port    int
	Logger  logr.Logger
}

// NewProvider creates a new Ollama provider
//
// TODO:
// - need to handle base URL better (with trailing slashes, etc.)
// - need to construct actual URL using baseURL, port, etc.
func NewProvider(opts *ProviderOpts) *Provider {
	opts.Logger.V(0).Info("Creating new OllamaProvider")

	client := client.NewClient(
		client.WithBaseURL("http://localhost:11434/api"),
	)

	return &Provider{
		host:   opts.BaseURL,
		port:   opts.Port,
		client: client,
		logger: opts.Logger,
	}
}

func (p *Provider) GetCapabilities(ctx context.Context) (*core.Capabilities, error) {
	p.logger.V(1).Info("Fetching capabilities")
	// Placeholder for future implementation
	p.logger.V(1).Info("GetCapabilities method is not implemented yet")
	return nil, nil
}

func (p *Provider) UseModel(ctx context.Context, model *model.Model) error {
	p.logger.V(1).Info("Setting model", "modelID", model.ID)
	p.model = model
	return nil
}

// Generate implements the LLMProvider interface for basic responses
func (p *Provider) Generate(ctx context.Context, opts *core.GenerateOptions) (*message.Message, error) {
	fmt.Printf("Messages in generate: %v\n", opts.Messages)
	p.logger.V(1).Info("Generate request received", "modelID", p.model.ID)
	ollamaMessages := convertManyMessagesToOllamaMessages(opts.Messages)

	// Convert tools into Ollama's format
	ollamaTools := make([]client.Tool, len(opts.Tools))
	p.logger.V(2).Info("Converting tools to Ollama format", "toolCount", len(opts.Tools))

	for i, t := range opts.Tools {
		ollamaTools[i] = client.Tool{
			Type: "function",
			Function: client.ToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.JSONSchema,
			},
		}
	}

	resp, err := p.client.Chat(ctx, &client.ChatRequest{
		Model:    p.model.ID,
		Messages: ollamaMessages,
		Tools:    ollamaTools,
	})
	if err != nil {
		p.logger.Error(err, "Error calling client chat method")
		return nil, fmt.Errorf("error calling client chat method: %w", err)
	}

	toolCalls := []message.ToolCall{}
	for _, toolCall := range resp.Message.ToolCalls {
		toolCalls = append(toolCalls, message.ToolCall{
			Name:      toolCall.Function.Name,
			Arguments: toolCall.Function.Arguments,
		})
	}

	return &message.Message{
		Role:      message.AssistantMessageRole,
		Content:   resp.Message.Content,
		ToolCalls: toolCalls,
	}, nil
}

// GenerateStream streams the response token by token
func (p *Provider) GenerateStream(ctx context.Context, opts *core.GenerateOptions) (<-chan *message.Message, <-chan error) {
	p.logger.V(1).Info("Stream generation not implemented yet")
	return nil, nil
}
