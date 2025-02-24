package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/ai-flowx/drivex/pkg/mylog"
)

func SendHTTPRequest(apiKey, url string, reqBody []byte, httpTransport *http.Transport) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}

	if httpTransport != nil {
		client.Transport = httpTransport
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := string(respBody)
		return nil, fmt.Errorf("http status code: %d, %s", resp.StatusCode, errMsg)
	}

	return respBody, nil
}

func SendSSERequest(apiKey, url string, reqBody []byte, callback func(data string), httpTransport *http.Transport) error {
	mylog.Logger.Debug("SendSSERequest", zap.String("url", url))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	if httpTransport != nil {
		client.Transport = httpTransport
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errMsg string
		respBody, err := io.ReadAll(resp.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		if err != nil {
			mylog.Logger.Error(err.Error())
		}
		if len(respBody) > 0 {
			errMsg = string(respBody)
		} else {
			errMsg = "empty response body"
		}

		return fmt.Errorf("http status code: %d, %s", resp.StatusCode, errMsg)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		mylog.Logger.Debug("SendSSERequest", zap.String("line", line))
		if err != nil {
			break
		}
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(line[5:])
			callback(data)
		}
	}

	return nil
}

func SendSSERequestWithHttpHeader(apiKey, url string, reqBody []byte, callback func(data string),
	httpTransport *http.Transport, header map[string]string) error {
	mylog.Logger.Debug("SendSSERequest", zap.String("url", url))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	for k, v := range header {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	if httpTransport != nil {
		client.Transport = httpTransport
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(line[5:])
			callback(data)
		}
	}

	return nil
}
