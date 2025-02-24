package mycommon

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"

	"github.com/ai-flowx/drivex/pkg/mylog"
)

const (
	roleAssistant = "assistant"
	roleSystem    = "system"
	roleUser      = "user"
)

func GetSystemMessage(oaiReqMessage []openai.ChatCompletionMessage) string {
	for i := 0; i < len(oaiReqMessage); i++ {
		msg := oaiReqMessage[i]
		if msg.Role == openai.ChatMessageRoleSystem {
			if len(msg.MultiContent) > 0 {
				for j := 0; j < len(msg.MultiContent); j++ {
					if msg.MultiContent[j].Type == openai.ChatMessagePartTypeText {
						return msg.MultiContent[j].Text
					}
				}
			} else {
				return msg.Content
			}
		}
	}

	return ""
}

func GetLastestMessage(oaiReqMessage []openai.ChatCompletionMessage) string {
	if len(oaiReqMessage) == 0 {
		return ""
	}

	lastestMsg := oaiReqMessage[len(oaiReqMessage)-1]

	if lastestMsg.Role == openai.ChatMessageRoleSystem {
		return ""
	} else {
		if len(lastestMsg.MultiContent) > 0 {
			for j := 0; j < len(lastestMsg.MultiContent); j++ {
				if lastestMsg.MultiContent[j].Type == openai.ChatMessagePartTypeText {
					return lastestMsg.MultiContent[j].Text
				}
			}
		} else {
			return lastestMsg.Content
		}
	}

	return ""
}

func IsMultiContentMessage(oaiReqMessage []openai.ChatCompletionMessage) bool {
	if len(oaiReqMessage) > 0 {
		for i := 0; i < len(oaiReqMessage); i++ {
			if len(oaiReqMessage[i].MultiContent) > 0 {
				return true
			}
		}
	}

	return false
}

func ConvertSystemMessages2NoSystem(oaiReqMessage []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	var systemQuery string
	if len(oaiReqMessage) == 0 {
		return oaiReqMessage
	}

	if strings.EqualFold(oaiReqMessage[0].Role, roleSystem) {
		if len(oaiReqMessage) == 1 {
			oaiReqMessage[0].Role = roleUser
		} else {
			systemQuery = oaiReqMessage[0].Content
			oaiReqMessage = oaiReqMessage[1:]
			oaiReqMessage[0].Content = systemQuery + "\n" + oaiReqMessage[0].Content
		}
	}

	mylog.Logger.Debug("ConvertSystemMessages2NoSystem", zap.Any("oaiReqMessage", oaiReqMessage))

	return oaiReqMessage
}

func NormalizeMessages(oaiReqMessage []openai.ChatCompletionMessage, keepAllSystem bool) []openai.ChatCompletionMessage {
	if len(oaiReqMessage) == 0 {
		return oaiReqMessage
	}

	if strings.EqualFold(oaiReqMessage[0].Role, roleSystem) {
		if len(oaiReqMessage) == 1 {
			oaiReqMessage[0].Role = roleUser
		}
	}

	var normalizedMessages []openai.ChatCompletionMessage
	var lastRole string

	for i := range oaiReqMessage {
		role := strings.ToLower(oaiReqMessage[i].Role)
		if !keepAllSystem && role == roleSystem && i > 0 {
			continue
		}
		if role == roleUser || role == roleAssistant {
			if role == lastRole {
				continue
			}
			normalizedMessages = append(normalizedMessages, oaiReqMessage[i])
		} else {
			normalizedMessages = append(normalizedMessages, oaiReqMessage[i])
		}
		lastRole = role
	}

	return normalizedMessages
}

func GetImageURLData(dataStr string) (base64Data, mime string, err error) {
	if strings.HasPrefix(dataStr, "data:") {
		sepIndex := strings.Index(dataStr, ",")
		if sepIndex == -1 {
			return "", "", fmt.Errorf("invalid data URL format")
		}
		mime = dataStr[5:sepIndex]
		base64Data = dataStr[sepIndex+1:]
		return base64Data, mime, nil
	} else if strings.HasPrefix(dataStr, "http") {
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		response, err := client.Get(dataStr)
		if err != nil {
			return "", "", fmt.Errorf("error fetching image: %v", err)
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(response.Body)
		if response.StatusCode != http.StatusOK {
			return "", "", fmt.Errorf("failed to download image: HTTP status %d", response.StatusCode)
		}
		var base64Writer strings.Builder
		encoder := base64.NewEncoder(base64.StdEncoding, &base64Writer)
		defer func(encoder io.WriteCloser) {
			_ = encoder.Close()
		}(encoder)
		if _, err := io.Copy(encoder, response.Body); err != nil {
			return "", "", fmt.Errorf("error encoding image data to base64: %v", err)
		}
		mimeType := response.Header.Get("Content-Type")
		return base64Writer.String(), mimeType, nil
	}

	return "", "", fmt.Errorf("unsupported URL format")
}

func AdjustOpenAIRequestParams(oaiReq *openai.ChatCompletionRequest) {
	adjustedTemperature, adjustedTopP, adjustedMaxTokens, err := AdjustParamsToRange(oaiReq.Model, oaiReq.Temperature,
		oaiReq.TopP, oaiReq.MaxTokens)
	if err != nil {
		return
	}

	oaiReq.Temperature = adjustedTemperature
	oaiReq.TopP = adjustedTopP
	oaiReq.MaxTokens = adjustedMaxTokens

	mylog.Logger.Debug("", zap.Float32("adjustedTemperature", adjustedTemperature),
		zap.Float32("adjustedTopP", adjustedTopP),
		zap.Int("MaxTokens", adjustedMaxTokens),
	)
}

func DeepCopyChatCompletionRequest(r *openai.ChatCompletionRequest) openai.ChatCompletionRequest {
	newRequest := *r
	newRequest.Messages = make([]openai.ChatCompletionMessage, len(r.Messages))

	for i := range r.Messages {
		newRequest.Messages[i] = r.Messages[i]
		if len(newRequest.Messages[i].MultiContent) > 0 {
			newRequest.Messages[i].MultiContent = make([]openai.ChatMessagePart, len(r.Messages[i].MultiContent))
			for j, part := range r.Messages[i].MultiContent {
				newRequest.Messages[i].MultiContent[j] = part
				if part.ImageURL != nil {
					newImageURL := *part.ImageURL
					newRequest.Messages[i].MultiContent[j].ImageURL = &newImageURL
				}
			}
		}
	}

	return newRequest
}

func LogChatCompletionRequest(request *openai.ChatCompletionRequest) {
	mylog.Logger.Debug("LogChatCompletionRequest", zap.Any("req", request))

	filteredRequest := DeepCopyChatCompletionRequest(request)

	for i := range filteredRequest.Messages {
		if len(filteredRequest.Messages[i].MultiContent) > 0 {
			for j, part := range filteredRequest.Messages[i].MultiContent {
				if part.Type == openai.ChatMessagePartTypeImageURL && part.ImageURL != nil {
					if !strings.HasPrefix(part.ImageURL.URL, "http") {
						d := "..."
						filteredRequest.Messages[i].MultiContent[j].ImageURL.URL = d
					}
				}
			}
		}
	}

	mylog.Logger.Debug("LogChatCompletionRequest", zap.Any("filteredRequest", filteredRequest))

	jsonData, err := json.Marshal(filteredRequest)
	if err != nil {
		mylog.Logger.Error("LogChatCompletionRequest|Marshal", zap.Error(err))
		return
	}

	mylog.Logger.Info("LogChatCompletionRequest", zap.String("request", string(jsonData)))
}

func ParseChatCompletionRequest(data []byte) (*openai.ChatCompletionRequest, error) {
	var rawRequest struct {
		Model       string            `json:"model"`
		Messages    []json.RawMessage `json:"messages"`
		Temperature float32           `json:"temperature,omitempty"`
		Stream      bool              `json:"stream,omitempty"`
	}

	if err := json.Unmarshal(data, &rawRequest); err != nil {
		return nil, err
	}

	request := &openai.ChatCompletionRequest{
		Model:       rawRequest.Model,
		Temperature: rawRequest.Temperature,
		Stream:      rawRequest.Stream,
	}

	for _, rawMsg := range rawRequest.Messages {
		var rawMessage struct {
			Role    string          `json:"role"`
			Content json.RawMessage `json:"content"`
		}
		if err := json.Unmarshal(rawMsg, &rawMessage); err != nil {
			return nil, err
		}
		message := openai.ChatCompletionMessage{
			Role: rawMessage.Role,
		}
		var contentStr string
		if err := json.Unmarshal(rawMessage.Content, &contentStr); err == nil {
			message.Content = contentStr
		} else {
			var contentObj struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}
			if err := json.Unmarshal(rawMessage.Content, &contentObj); err == nil {
				if contentObj.Type == string(openai.ChatMessagePartTypeText) {
					message.Content = contentObj.Text
				} else {
					return nil, fmt.Errorf("unexpected content type: %s", contentObj.Type)
				}
			} else {
				var contentArr []openai.ChatMessagePart
				if err := json.Unmarshal(rawMessage.Content, &contentArr); err == nil {
					messageBytes, err := json.Marshal(contentArr)
					if err != nil {
						return nil, fmt.Errorf("failed to marshal content array")
					}
					message.Content = string(messageBytes)
				} else {
					return nil, fmt.Errorf("failed to unmarshal content")
				}
			}
		}

		request.Messages = append(request.Messages, message)
	}

	return request, nil
}
