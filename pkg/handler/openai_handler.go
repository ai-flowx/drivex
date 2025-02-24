package handler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"

	"github.com/ai-flowx/drivex/pkg/config"
	"github.com/ai-flowx/drivex/pkg/mycommon"
	"github.com/ai-flowx/drivex/pkg/mylimiter"
	"github.com/ai-flowx/drivex/pkg/mylog"
	"github.com/ai-flowx/drivex/pkg/utils"
)

const (
	defaultReqTimeout = 10
)

type OAIRequestParam struct {
	chatCompletionReq *openai.ChatCompletionRequest
	modelDetails      *config.ModelDetails
	creds             map[string]interface{}
	httpTransport     *http.Transport
	ClientModel       string
}

var serviceHandlerMap = map[string]func(*gin.Context, *OAIRequestParam) error{
	"siliconflow": OpenAI2SiliconFlowHandler,
}

func LogRequestDetails(c *gin.Context) {
	mylog.Logger.Debug("HTTP request details",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Any("parameters", c.Request.URL.Query()),
		zap.Any("headers", c.Request.Header),
	)
}

func logOpenAIChatCompletionRequest(oaiReq *openai.ChatCompletionRequest) {
	if oaiReq == nil {
		return
	}

	mylog.Logger.Info("logOpenAIChatCompletionRequest", zap.Float32("TopP", oaiReq.TopP),
		zap.Float32("Temperature", oaiReq.Temperature), zap.Int("MaxTokens", oaiReq.MaxTokens),
		zap.String("model", oaiReq.Model), zap.Int("N", oaiReq.N), zap.Float32("FrequencyPenalty", oaiReq.FrequencyPenalty))
}

func getBodyDataCopy(c *gin.Context) ([]byte, error) {
	body, err := c.GetRawData()
	if err != nil {
		return nil, err
	}

	c.Set("rawData", body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

func OpenAIHandler(c *gin.Context) {
	if !validateRequestMethod(c, "POST") {
		return
	}

	LogRequestDetails(c)

	apikey, err := utils.GetAPIKeyFromHeader(c)
	if err != nil {
		mylog.Logger.Error(err.Error())
	}

	mylog.Logger.Info("OpenAIHandler", zap.String("apikey", apikey))

	isValid := validateAPIKey(apikey)
	if !isValid {
		err = errors.New("key is not valid")
		mylog.Logger.Error("key is not valid", zap.String("apikey", apikey))
		sendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bodyData, getBodyerr := getBodyDataCopy(c)

	var oaiReq openai.ChatCompletionRequest

	if err := c.ShouldBindJSON(&oaiReq); err != nil {
		mylog.Logger.Error(err.Error())
		if getBodyerr != nil {
			mylog.Logger.Error(err.Error())
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		mylog.Logger.Debug(string(bodyData))
		parsedReq, parseErr := mycommon.ParseChatCompletionRequest(bodyData)
		if parseErr != nil {
			mylog.Logger.Error("ParseChatCompletionRequest error: " + parseErr.Error())
			sendErrorResponse(c, http.StatusBadRequest, parseErr.Error())
			return
		}
		oaiReq = *parsedReq
	}

	mylog.Logger.Info("logOpenAIChatCompletionRequest", zap.Float32("TopP", oaiReq.TopP))
	logOpenAIChatCompletionRequest(&oaiReq)

	if isValid, _ = config.ValidateAPIKeyAndModel(apikey, oaiReq.Model); !isValid {
		err = errors.New("key not valid")
		mylog.Logger.Error(err.Error())
		sendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	mycommon.LogChatCompletionRequest(&oaiReq)

	HandleOpenAIRequest(c, &oaiReq)
}

// nolint:funlen,gocyclo
func HandleOpenAIRequest(c *gin.Context, oaiReq *openai.ChatCompletionRequest) {
	clientModel := oaiReq.Model
	gRedirectModel := config.GetGlobalModelRedirect(clientModel)
	oaiReq.Model = gRedirectModel

	s, serviceModelName, err := getModelDetails(oaiReq)
	if err != nil {
		mylog.Logger.Error(err.Error())
		sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	mrModel := config.GetModelRedirect(s, serviceModelName)
	mpModel := config.GetModelMapping(s, mrModel)
	oaiReq.Model = mpModel

	mylog.Logger.Info("Service details",
		zap.String("service_name", s.ServiceName),
		zap.String("client_model", clientModel),
		zap.String("g_redirect_model", gRedirectModel),
		zap.String("service_model_name", serviceModelName),
		zap.String("redirect_model", mrModel),
		zap.String("map_model", mpModel),
		zap.String("last_model", oaiReq.Model))

	if mycommon.IsMultiContentMessage(oaiReq.Messages) {
		// TBD: FIXME
		mylog.Logger.Error("Not support multi content message")
		return
	}

	creds, credsID := mycommon.GetACredentials(s)

	var limiter *mylimiter.Limiter

	lt, ln, timeout := mycommon.GetServiceModelDetailsLimit(s)
	if lt != "" && ln > 0 {
		limiter = mylimiter.GetLimiter(s.ServiceID, lt, ln)
	} else {
		lt, ln, timeout = mycommon.GetCredentialLimit(creds)
		if lt != "" && ln > 0 {
			limiter = mylimiter.GetLimiter(credsID, lt, ln)
		}
	}

	oaiReqParam := &OAIRequestParam{
		chatCompletionReq: oaiReq,
		modelDetails:      s,
		creds:             creds,
		ClientModel:       clientModel,
	}

	if limiter != nil {
		if timeout <= 0 {
			timeout = defaultReqTimeout
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		startWaitTime := time.Now()
		mylog.Logger.Info("Rate limits and timeout configuration",
			zap.String("limit type:", lt),
			zap.Float64("limit num:", ln),
			zap.Int("timeout", timeout))
		if lt == "qps" || lt == "qpm" {
			err = limiter.Wait(ctx)
			elapsed := time.Since(startWaitTime)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					mylog.Logger.Error("Failed to obtain token within the specified time",
						zap.Error(err),
						zap.Int("timeout", timeout),
						zap.Duration("elapsed", elapsed))
				} else if errors.Is(err, context.Canceled) {
					mylog.Logger.Error("Operation canceled %v, actual waiting time: %v", zap.Error(err), zap.Duration("elapsed", elapsed))
				} else {
					mylog.Logger.Error("Unknown error occurred while waiting for a token: ", zap.Error(err), zap.Duration("elapsed", elapsed))
				}
				mylog.Logger.Info("waited for: ", zap.Duration("elapsed", elapsed))
				sendErrorResponse(c, http.StatusTooManyRequests, "Request rate limit exceeded")
				return
			}
			mylog.Logger.Info("Wait duration",
				zap.Duration("waited_for", time.Since(startWaitTime)))
		} else if lt == "concurrency" {
			err := limiter.Acquire(ctx)
			if err != nil {
				mylog.Logger.Error(err.Error())
			}
			defer limiter.Release()
			mylog.Logger.Info("Concurrency wait time",
				zap.Duration("waited_for", time.Since(startWaitTime)))
		}
	}

	if config.IsProxyEnabled(s) {
		proxyType, proxyAddr, transport, err := config.GetConfProxyTransport()
		if err != nil {
			mylog.Logger.Error("GetConfProxyTransport", zap.Error(err))
		} else {
			mylog.Logger.Debug("GetConfProxyTransport", zap.String("proxyType", proxyType), zap.String("proxyAddr", proxyAddr))
			oaiReqParam.httpTransport = transport
		}
	} else {
		mylog.Logger.Debug("GetConfProxyTransport proxy not enabled")
	}

	oaiReq.Messages = mycommon.NormalizeMessages(oaiReq.Messages, false)

	if err := dispatchToServiceHandler(c, oaiReqParam); err != nil {
		mylog.Logger.Error(err.Error())
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if oaiReq.Stream {
		utils.SendOpenAIStreamEOFData(c)
	}
}

func dispatchToServiceHandler(c *gin.Context, oaiReqParam *OAIRequestParam) error {
	s := oaiReqParam.modelDetails
	serviceName := strings.ToLower(s.ServiceName)

	if handler, ok := serviceHandlerMap[serviceName]; ok {
		return handler(c, oaiReqParam)
	}

	return errors.New("service handler not found")
}

func validateRequestMethod(c *gin.Context, method string) bool {
	if c.Request.Method != method {
		sendErrorResponse(c, http.StatusMethodNotAllowed, "Only "+method+" method is accepted")
		return false
	}

	return true
}

func validateAPIKey(apikey string) bool {
	if config.APIKey == "" {
		return true
	}

	if config.APIKey != apikey {
		return false
	}

	return true
}

func getModelDetails(oaiReq *openai.ChatCompletionRequest) (*config.ModelDetails, string, error) {
	if oaiReq.Model == config.KeynameRandom {
		return config.GetRandomEnabledModelDetailsV1()
	}

	s, err := config.GetModelService(oaiReq.Model)
	if err != nil {
		return nil, "", err
	}

	return s, oaiReq.Model, err
}

func sendErrorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}
