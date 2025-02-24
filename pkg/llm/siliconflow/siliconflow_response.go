package siliconflow

type ChatCompletionResponse struct {
	ID        string             `json:"id"`
	Choices   []ChoiceResponse   `json:"choices"`
	ToolCalls []ToolCallResponse `json:"tool_calls"`
	Usage     UsageResponse      `json:"usage"`
	Created   int64              `json:"created"`
	Model     string             `json:"model"`
	Object    string             `json:"object"`
}

type ChoiceResponse struct {
	Message      MessageResponse `json:"message,omitempty"`
	FinishReason string          `json:"finish_reason,omitempty"`
}

type MessageResponse struct {
	Role             string `json:"role,omitempty"`
	Content          string `json:"content,omitempty"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

type ToolCallResponse struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function FunctionResponse `json:"function"`
}

type FunctionResponse struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type UsageResponse struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}
