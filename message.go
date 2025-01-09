package ollama

import "github.com/agent-api/core/message"

type OllamaMessageRole string

const (
	OllamaUserMessageRole OllamaMessageRole = "user"
	OllamaAIMessageRole   OllamaMessageRole = "ai"
)

type ollamaMessage struct {
	Role    OllamaMessageRole `json:"role"`
	Content string            `json:"content"`
}

func convertMessageToOllamaMessage(m message.Message) *ollamaMessage {
	if m.Role == message.UserMessageRole {
		return &ollamaMessage{
			Role:    OllamaUserMessageRole,
			Content: m.Content,
		}
	}

	return nil
}

func convertOllamaMessageToMessage(m ollamaMessage) message.Message {
	if m.Role == OllamaUserMessageRole {
		return message.Message{
			Role:    message.UserMessageRole,
			Content: m.Content,
		}
	}

	return message.Message{}
}

func convertManyMessagesToOllamaMessages(messages []message.Message) []*ollamaMessage {
	// Convert agent messages to Ollama format
	ollamaMessages := make([]*ollamaMessage, len(messages))

	for i, m := range messages {
		ollamaMessages[i] = convertMessageToOllamaMessage(m)
	}

	return ollamaMessages
}
