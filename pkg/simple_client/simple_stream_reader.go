package simple_client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	writer io.Writer
	status int
	header http.Header
	body   *bytes.Buffer
}

func NewCustomResponseWriter(w io.Writer) *CustomResponseWriter {
	return &CustomResponseWriter{
		writer: w,
		header: http.Header{},
		body:   bytes.NewBuffer([]byte{}),
	}
}

func (crw *CustomResponseWriter) CloseNotify() <-chan bool {
	if notifier, ok := crw.writer.(http.CloseNotifier); ok {
		return notifier.CloseNotify()
	}

	c := make(chan bool)
	close(c)

	return c
}

func (crw *CustomResponseWriter) Write(data []byte) (int, error) {
	crw.body.Write(data)

	return crw.writer.Write(data)
}

func (crw *CustomResponseWriter) WriteHeader(statusCode int) {
	crw.status = statusCode

	_, _ = fmt.Fprintf(crw.writer, "HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))
}

func (crw *CustomResponseWriter) WriteString(s string) (int, error) {
	return crw.WriteString(s)
}

func (crw *CustomResponseWriter) Header() http.Header {
	return http.Header{}
}

func (crw *CustomResponseWriter) Status() int {
	return crw.status
}

func (crw *CustomResponseWriter) Size() int {
	return crw.body.Len()
}

func (crw *CustomResponseWriter) Flush() {
	if flusher, ok := crw.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

type SimpleChatCompletionStream struct {
	reader *bufio.Reader
}

func NewSimpleChatCompletionStream(reader io.Reader) *SimpleChatCompletionStream {
	return &SimpleChatCompletionStream{reader: bufio.NewReader(reader)}
}

func (scs *SimpleChatCompletionStream) Recv() (*openai.ChatCompletionStreamResponse, error) {
	var response openai.ChatCompletionStreamResponse

	line, err := scs.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, err
	}

	if len(line) == 1 && string(line) == "\n" {
		return nil, nil
	}

	if strings.Contains(string(line), "[DONE]") {
		return nil, io.EOF
	}

	data := strings.TrimSpace(string(line))
	if strings.HasPrefix(data, "data: ") {
		jsonData := strings.TrimPrefix(data, "data: ")
		if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
			return &response, err
		}
		return &response, nil
	}

	errData, _ := io.ReadAll(scs.reader)

	return &response, fmt.Errorf("unexpected data format: %s", string(errData))
}
