package adapter

import (
	"github.com/sashabaranov/go-openai"

	"github.com/ai-flowx/drivex/pkg/llm/siliconflow"
	myopenai "github.com/ai-flowx/drivex/pkg/openai"
)

func OpenAIRequestToSiliconFlowRequest(req *openai.ChatCompletionRequest) *siliconflow.ChatCompletionRequest {
	return &siliconflow.ChatCompletionRequest{}
}

func SiliconFlowReponseToOpenAIResponse(resp *siliconflow.ChatCompletionResponse) *myopenai.OpenAIResponse {
	return &myopenai.OpenAIResponse{}
}
