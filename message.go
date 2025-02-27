package ollama

import (
	"github.com/agent-api/core/types"
	"github.com/agent-api/ollama/client"
)

func convertMessageToOllamaMessage(m *types.Message) *client.Message {
	images := []string{}
	for _, image := range m.Images {
		images = append(images, image.Base64Encoding)
	}

	switch m.Role {
	case types.UserMessageRole:
		return &client.Message{
			Role:    client.RoleUser,
			Content: m.Content,
			Images:  images,
		}

	case types.AssistantMessageRole:
		return &client.Message{
			Role:    client.RoleAssistant,
			Content: m.Content,
			Images:  images,
		}

	case types.ToolMessageRole:
		return &client.Message{
			Role:    client.RoleTool,
			Content: m.Content,
			Images:  images,
		}
	}

	return nil
}

func convertOllamaMessageToMessage(m *client.Message) *types.Message {
	images := []*types.Image{}
	for _, image := range m.Images {
		images = append(images, &types.Image{
			Base64Encoding: image,
		})
	}

	switch m.Role {
	case client.RoleUser:
		return &types.Message{
			Role:    types.UserMessageRole,
			Content: m.Content,
			Images:  images,
		}

	case client.RoleAssistant:
		return &types.Message{
			Role:    types.AssistantMessageRole,
			Content: m.Content,
			Images:  images,
		}

	case client.RoleTool:
		return &types.Message{
			Role:    types.ToolMessageRole,
			Content: m.Content,
			Images:  images,
		}
	}

	return nil
}

func convertManyMessagesToOllamaMessages(messages []*types.Message) []*client.Message {
	// Convert agent messages to Ollama format
	ollamaMessages := make([]*client.Message, len(messages))

	for i, m := range messages {
		ollamaMessages[i] = convertMessageToOllamaMessage(m)
	}

	return ollamaMessages
}
