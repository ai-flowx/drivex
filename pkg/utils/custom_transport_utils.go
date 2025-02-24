package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CustomTransport struct {
	Transport http.RoundTripper
}

func (c *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("error reading error response body: %v", readErr)
		}
		_ = resp.Body.Close()
		return nil, fmt.Errorf("HTTP error: %s, body: %s", resp.Status, string(bodyBytes))
	}

	modifiedBody := &modifiedReadCloser{
		originalBody: resp.Body,
		reader:       bufio.NewReader(resp.Body),
	}

	resp.Body = modifiedBody

	return resp, nil
}

type modifiedReadCloser struct {
	originalBody io.ReadCloser
	buf          *bytes.Buffer
	reader       *bufio.Reader
}

func (m *modifiedReadCloser) Read(p []byte) (int, error) {
	if m.buf == nil || m.buf.Len() == 0 {
		line, err := m.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}
			return 0, err
		}
		if strings.HasPrefix(line, "data:") && !strings.HasPrefix(line, "data: ") {
			line = strings.Replace(line, "data:", "data: ", 1)
		}
		m.buf = bytes.NewBufferString(line)
	}

	return m.buf.Read(p)
}

func (m *modifiedReadCloser) Close() error {
	return m.originalBody.Close()
}
