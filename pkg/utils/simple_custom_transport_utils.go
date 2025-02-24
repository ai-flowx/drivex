package utils

import (
	"fmt"
	"io"
	"net/http"
)

const (
	bodyLen = 1024
)

type SimpleCustomTransport struct {
	Transport http.RoundTripper
}

func (c *SimpleCustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		bodyBytes := make([]byte, bodyLen)
		n, readErr := resp.Body.Read(bodyBytes)
		if readErr != nil && readErr != io.EOF {
			return nil, fmt.Errorf("error reading error response body: %v", readErr)
		}
		_ = resp.Body.Close()
		return nil, fmt.Errorf("HTTP error: %s, body: %s", resp.Status, string(bodyBytes[:n]))
	}

	return resp, nil
}
