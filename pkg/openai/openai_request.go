package myopenai

import (
	"encoding/json"

	"github.com/ai-flowx/drivex/pkg/mycommon"
)

type OpenAIRequest struct {
	Model            string             `json:"model"`
	Messages         []mycommon.Message `json:"messages"`
	FrequencyPenalty *float32           `json:"frequency_penalty,omitempty"`
	LogitBias        map[int]int        `json:"logit_bias,omitempty"`
	LogProbs         *bool              `json:"logprobs,omitempty"`
	TopLogProbs      *int               `json:"top_logprobs,omitempty"`
	MaxTokens        *int               `json:"max_tokens,omitempty"`
	N                *int               `json:"n,omitempty"`
	PresencePenalty  *float32           `json:"presence_penalty,omitempty"`
	ResponseFormat   *ResponseFormat    `json:"response_format,omitempty"`
	Seed             *int               `json:"seed,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	Stream           *bool              `json:"stream,omitempty"`
	StreamOptions    *StreamOptions     `json:"stream_options,omitempty"`
	Temperature      *float32           `json:"temperature,omitempty"`
	TopP             *float32           `json:"top_p,omitempty"`
	Tools            []Tool             `json:"tools,omitempty"`
	ToolChoice       json.RawMessage    `json:"tool_choice,omitempty"`
	User             *string            `json:"user,omitempty"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type StreamOptions struct {
	// 详细字段可以根据具体实现需求添加
}

type Tool struct {
	Type     string    `json:"type"`
	Function *Function `json:"function,omitempty"`
}

type Function struct {
	Name string `json:"name"`
}

type ToolChoiceFunction struct {
	Type     string    `json:"type"`
	Function *Function `json:"function,omitempty"`
}
