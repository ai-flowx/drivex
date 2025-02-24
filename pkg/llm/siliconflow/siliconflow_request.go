package siliconflow

type ChatCompletionRequest struct {
	Model            string                `json:"model"`
	Messages         []MessageRequest      `json:"messages"`
	Stream           bool                  `json:"stream,omitempty"`
	MaxTokens        int                   `json:"max_tokens,omitempty"`
	Stop             []string              `json:"stop,omitempty"`
	Temperature      float32               `json:"temperature,omitempty"`
	TopP             float32               `json:"top_p,omitempty"`
	TopK             float32               `json:"top_k,omitempty"`
	FrequencyPenalty float32               `json:"frequency_penalty,omitempty"`
	N                int                   `json:"n,omitempty"`
	ResponseFormat   ResponseFormatRequest `json:"response_format,omitempty"`
	Tools            []ToolRequest         `json:"tools,omitempty"`
}

type MessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormatRequest struct {
	Type string `json:"type,omitempty"`
}

type ToolRequest struct {
	Type     string          `json:"type"`
	Function FunctionRequest `json:"function"`
}

type FunctionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
	Strict      bool   `json:"strict,omitempty"`
}
