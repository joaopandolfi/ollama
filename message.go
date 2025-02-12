package ollama

import (
	"github.com/agent-api/core/message"
	"github.com/agent-api/ollama-provider/client"
)

func convertMessageToOllamaMessage(m message.Message) *client.Message {
	switch m.Role {
	case message.UserMessageRole:
		return &client.Message{
			Role:    client.RoleUser,
			Content: m.Content,
		}

	case message.AssistantMessageRole:
		return &client.Message{
			Role:    client.RoleAssistant,
			Content: m.Content,
		}

	case message.ToolMessageRole:
		return &client.Message{
			Role:    client.RoleTool,
			Content: m.Content,
		}
	}

	return nil
}

func convertOllamaMessageToMessage(m client.Message) message.Message {
	switch m.Role {
	case client.RoleUser:
		return message.Message{
			Role:    message.UserMessageRole,
			Content: m.Content,
		}

	case client.RoleAssistant:
		return message.Message{
			Role:    message.AssistantMessageRole,
			Content: m.Content,
		}

	case client.RoleTool:
		return message.Message{
			Role:    message.ToolMessageRole,
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
