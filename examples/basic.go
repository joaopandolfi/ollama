package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/agent-api/ollama-provider"
)

func main() {
	// Create a new client
	client := ollama.NewClient(
		ollama.WithBaseURL("http://localhost:11434/api"),
	)

	// Define a tool
	weatherTool := ollama.Tool{
		Type: "function",
		Function: ollama.ToolFunction{
			Name:        "get_weather",
			Description: "Get the current weather for a location",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"location": {
						"type": "string",
						"description": "The city and state/country"
					}
				},
				"required": ["location"]
			}`),
		},
	}

	// Create messages
	messages := []ollama.Message{
		{
			Role:    ollama.RoleUser,
			Content: "What's the weather like in San Francisco?",
		},
	}

	// Stream the chat responses
	_, err := client.ChatStream(
		context.Background(),
		ollama.ChatRequest{
			Model:    "llama2",
			Messages: messages,
			Tools:    []ollama.Tool{weatherTool},
		},
		func(response *ollama.ChatResponse) error {
			// Handle tool calls if present
			if len(response.Message.ToolCalls) > 0 {
				for _, call := range response.Message.ToolCalls {
					fmt.Printf("Tool called: %s with args: %s\n",
						call.Function.Name,
						string(call.Function.Arguments),
					)
				}
			}

			// Print the response content
			fmt.Print(response.Message.Content)
			return nil
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}
