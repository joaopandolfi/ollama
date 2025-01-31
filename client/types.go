package client

import (
	"time"
)

// RequestOptions represents optional parameters for the request
type RequestOptions struct {
	NumKeep          *int     `json:"num_keep,omitempty"`
	Seed             *int     `json:"seed,omitempty"`
	NumPredict       *int     `json:"num_predict,omitempty"`
	TopK             *int     `json:"top_k,omitempty"`
	TopP             *float64 `json:"top_p,omitempty"`
	MinP             *float64 `json:"min_p,omitempty"`
	TFSZ             *float64 `json:"tfs_z,omitempty"`
	TypicalP         *float64 `json:"typical_p,omitempty"`
	RepeatLastN      *int     `json:"repeat_last_n,omitempty"`
	Temperature      *float64 `json:"temperature,omitempty"`
	RepeatPenalty    *float64 `json:"repeat_penalty,omitempty"`
	PresencePenalty  *float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	Mirostat         *int     `json:"mirostat,omitempty"`
	MirostatTau      *float64 `json:"mirostat_tau,omitempty"`
	MirostatEta      *float64 `json:"mirostat_eta,omitempty"`
	PenalizeNewline  *bool    `json:"penalize_newline,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	Numa             *bool    `json:"numa,omitempty"`
	NumCtx           *int     `json:"num_ctx,omitempty"`
	NumBatch         *int     `json:"num_batch,omitempty"`
	NumGPU           *int     `json:"num_gpu,omitempty"`
	MainGPU          *int     `json:"main_gpu,omitempty"`
	LowVRAM          *bool    `json:"low_vram,omitempty"`
	F16KV            *bool    `json:"f16_kv,omitempty"`
	LogitsAll        *bool    `json:"logits_all,omitempty"`
	VocabOnly        *bool    `json:"vocab_only,omitempty"`
	UseMMap          *bool    `json:"use_mmap,omitempty"`
	UseMLock         *bool    `json:"use_mlock,omitempty"`
	NumThread        *int     `json:"num_thread,omitempty"`
}

// ChatRequest represents a request to the chat endpoint
type ChatRequest struct {
	Model     string          `json:"model"`
	Messages  []*Message      `json:"messages"`
	Format    *string         `json:"format,omitempty"`
	Options   *RequestOptions `json:"options,omitempty"`
	Stream    *bool           `json:"stream,omitempty"`
	KeepAlive *int            `json:"keep_alive,omitempty"`
	Tools     []Tool          `json:"tools,omitempty"`
}

// ChatResponse represents a response from the chat endpoint
type ChatResponse struct {
	Message         Message   `json:"message"`
	Model           string    `json:"model"`
	CreatedAt       time.Time `json:"created_at"`
	Done            bool      `json:"done"`
	DoneReason      string    `json:"done_reason,omitempty"`
	TotalDuration   int64     `json:"total_duration"`
	LoadDuration    int64     `json:"load_duration"`
	PromptEvalCount int       `json:"prompt_eval_count"`
	EvalCount       int       `json:"eval_count"`
	EvalDuration    int64     `json:"eval_duration"`
}
