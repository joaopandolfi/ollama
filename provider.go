package ollama

import (
	"context"
	"fmt"

	"github.com/agent-api/core"
	"github.com/agent-api/ollama/client"
	"github.com/go-logr/logr"
)

// Provider implements the LLMProvider interface for Ollama
type Provider struct {
	host string
	port int

	model *core.Model

	// client is the internal Ollama HTTP client
	client *client.OllamaClient

	logger *logr.Logger
}

type ProviderOpts struct {
	BaseURL string
	Port    int
	Logger  *logr.Logger
}

// NewProvider creates a new Ollama provider
//
// TODO:
// - need to handle base URL better (with trailing slashes, etc.)
// - need to construct actual URL using baseURL, port, etc.
func NewProvider(opts *ProviderOpts) *Provider {
	opts.Logger.Info("Creating new Provider")

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
	p.logger.Info("Fetching capabilities")

	// Placeholder for future implementation
	p.logger.Info("GetCapabilities method is not implemented yet")
	return nil, nil
}

func (p *Provider) UseModel(ctx context.Context, model *core.Model) error {
	p.logger.Info("Setting model", "modelID", model.ID)
	p.model = model
	return nil
}

// Generate implements the LLMProvider interface for basic responses
func (p *Provider) Generate(ctx context.Context, opts *core.GenerateOptions) (*core.Message, error) {
	p.logger.Info("Generate request received", "modelID", p.model.ID)
	ollamaMessages := convertManyMessagesToOllamaMessages(opts.Messages)

	// Convert tools into Ollama's format
	ollamaTools := make([]*client.Tool, len(opts.Tools))
	p.logger.Info("Converting tools to Ollama format", "toolCount", len(opts.Tools))

	for i, t := range opts.Tools {
		ollamaTools[i] = &client.Tool{
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
		p.logger.V(-1).Info(err.Error(), "client error", err)
		return nil, fmt.Errorf("error calling client chat method: %w", err)
	}

	toolCalls := []*core.ToolCall{}
	for _, toolCall := range resp.Message.ToolCalls {
		toolCalls = append(toolCalls, &core.ToolCall{
			Name:      toolCall.Function.Name,
			Arguments: toolCall.Function.Arguments,
		})
	}

	return &core.Message{
		Role:      core.AssistantMessageRole,
		Content:   resp.Message.Content,
		ToolCalls: toolCalls,
	}, nil
}

// GenerateStream streams the response token by token
func (p *Provider) GenerateStream(ctx context.Context, opts *core.GenerateOptions) (<-chan *core.Message, <-chan string, <-chan error) {
	p.logger.Info("Stream generation not implemented yet")
	return nil, nil, nil
}
