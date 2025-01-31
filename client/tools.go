package client

import "encoding/json"

// Tool represents a function tool that the model can use
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents a function that can be called by the model
type ToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

// ToolCall represents a tool call made by the model
type ToolCall struct {
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction represents the function being called
type ToolCallFunction struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}
