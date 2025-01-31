package ollama

import (
	"github.com/agent-api/core/message"
	"github.com/agent-api/ollama-provider/client"
)

func convertMessageToOllamaMessage(m message.Message) *client.Message {
	if m.Role == message.UserMessageRole {
		return &client.Message{
			Role:    client.RoleUser,
			Content: m.Content,
		}
	}

	return nil
}

func convertOllamaMessageToMessage(m client.Message) message.Message {
	if m.Role == client.RoleUser {
		return message.Message{
			Role:    message.UserMessageRole,
			Content: m.Content,
		}
	}

	return message.Message{}
}

func convertManyMessagesToOllamaMessages(messages []message.Message) []*client.Message {
	// Convert agent messages to Ollama format
	ollamaMessages := make([]*client.Message, len(messages))

	for i, m := range messages {
		ollamaMessages[i] = convertMessageToOllamaMessage(m)
	}

	return ollamaMessages
}
