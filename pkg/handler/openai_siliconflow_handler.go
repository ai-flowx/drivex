package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"

	myopenai "github.com/ai-flowx/drivex/pkg/openai"
)

var (
	defaultSiliconFlowServerURL = "https://api.siliconflow.cn/v1"
)

func OpenAI2SiliconFlowHandler(c *gin.Context, oaiReqParam *OAIRequestParam) error {
	return nil
}

func sendSiliconFlowRequest(c *gin.Context, client *http.Client, apiKey, url string, request interface{},
	oaiReq *openai.ChatCompletionRequest, oaiReqParam *OAIRequestParam) error {
	resp := http.Response{}
	return handleSiliconFlowResponse(c, &resp, oaiReq, oaiReqParam)
}

func handleSiliconFlowResponse(c *gin.Context, resp *http.Response, oaiReq *openai.ChatCompletionRequest,
	oaiReqParam *OAIRequestParam) error {
	return nil
}

func handleSiliconFlowStreamResponse(c *gin.Context, resp *http.Response, oaiReq *openai.ChatCompletionRequest,
	oaiReqParam *OAIRequestParam) error {
	// TBD: FIXME
	return nil
}

func processSiliconFlowStreamEvent(c *gin.Context, eventType, eventData, clientModel string) error {
	// TBD: FIXME
	return nil
}

func handleSiliconFlowEvent[T any](c *gin.Context, eventData string, eventStruct T,
	converter func(*T) *myopenai.OpenAIStreamResponse, clientModel string) error {
	// TBD: FIXME
	return nil
}
