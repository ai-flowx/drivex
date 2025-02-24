package simple_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"

	"github.com/ai-flowx/drivex/pkg/handler"
)

type SimpleClient struct {
}

func NewSimpleClient(authToken string) *SimpleClient {
	return NewSimpleClientWithConfig()
}

func NewSimpleClientWithConfig() *SimpleClient {
	return &SimpleClient{}
}

func (c *SimpleClient) CreateChatCompletion(
	_ context.Context,
	request *openai.ChatCompletionRequest,
) (response openai.ChatCompletionResponse, err error) {
	request.Stream = false
	reqBody, _ := json.Marshal(request)
	httpReq, _ := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	ginc := gin.New()
	ginc.POST("/v1/chat/completions", func(ctx *gin.Context) {
		handler.HandleOpenAIRequest(ctx, request)
	})

	w := httptest.NewRecorder()

	ginc.ServeHTTP(w, httpReq)

	if w.Code >= http.StatusBadRequest {
		err = errors.New(w.Body.String())
		return
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)

	return
}

func (c *SimpleClient) CreateChatCompletionStream(
	_ context.Context,
	request *openai.ChatCompletionRequest,
) (stream *SimpleChatCompletionStream, err error) {
	request.Stream = true

	reader, writer := io.Pipe()
	recorder := httptest.NewRecorder()

	ginc := gin.New()
	ginc.Use(func(ctx *gin.Context) {
		crw := NewCustomResponseWriter(writer)
		ctx.Writer = crw
		ctx.Next()
	})
	ginc.POST("/v1/chat/completions", func(ctx *gin.Context) {
		handler.HandleOpenAIRequest(ctx, request)
	})

	go func() {
		defer func(writer *io.PipeWriter) {
			_ = writer.Close()
		}(writer)
		requestData, _ := json.Marshal(request)
		httpReq, _ := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer(requestData))
		httpReq.Header.Set("Content-Type", "application/json")
		ginc.ServeHTTP(recorder, httpReq)
	}()

	return NewSimpleChatCompletionStream(reader), nil
}
