package ollama

import (
	"github.com/agent-api/core"
	"github.com/agent-api/ollama/client"
)

func convertMessageToOllamaMessage(m *core.Message) *client.Message {
	images := []string{}
	for _, image := range m.Images {
		images = append(images, image.Base64Encoding)
	}

	switch m.Role {
	case core.UserMessageRole:
		return &client.Message{
			Role:    client.RoleUser,
			Content: m.Content,
			Images:  images,
		}

	case core.AssistantMessageRole:
		return &client.Message{
			Role:    client.RoleAssistant,
			Content: m.Content,
			Images:  images,
		}

	case core.ToolMessageRole:
		return &client.Message{
			Role:    client.RoleTool,
			Content: m.Content,
			Images:  images,
		}
	}

	return nil
}

func convertOllamaMessageToMessage(m *client.Message) *core.Message {
	images := []*core.Image{}
	for _, image := range m.Images {
		images = append(images, &core.Image{
			Base64Encoding: image,
		})
	}

	switch m.Role {
	case client.RoleUser:
		return &core.Message{
			Role:    core.UserMessageRole,
			Content: m.Content,
			Images:  images,
		}

	case client.RoleAssistant:
		return &core.Message{
			Role:    core.AssistantMessageRole,
			Content: m.Content,
			Images:  images,
		}

	case client.RoleTool:
		return &core.Message{
			Role:    core.ToolMessageRole,
			Content: m.Content,
			Images:  images,
		}
	}

	return nil
}

func convertManyMessagesToOllamaMessages(messages []*core.Message) []*client.Message {
	// Convert agent messages to Ollama format
	ollamaMessages := make([]*client.Message, len(messages))

	for i, m := range messages {
		ollamaMessages[i] = convertMessageToOllamaMessage(m)
	}

	return ollamaMessages
}
