package deepseek

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type StreamChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content,omitempty"`
}

type StreamChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
	Stream   bool                    `json:"stream,omitempty"`
}

type ChatCompletionStream interface {
	Recv() (*StreamChatCompletionResponse, error)
	Close() error
}

// chatCompletionStream implements the ChatCompletionStream interface.
type chatCompletionStream struct {
	ctx    context.Context
	cancel context.CancelFunc
	resp   *http.Response
	reader *bufio.Reader
}

type StreamUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type StreamDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content"`
}

type StreamChoices struct {
	Index        int `json:"index"`
	Delta        StreamDelta
	FinishReason string `json:"finish_reason"`
}

type StreamChatCompletionResponse struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created int64           `json:"created"`
	Model   string          `json:"model"`
	Choices []StreamChoices `json:"choices"`
	Usage   *StreamUsage    `json:"usage,omitempty"`
}

func (s *chatCompletionStream) Recv() (*StreamChatCompletionResponse, error) {
	reader := bufio.NewReader(s.reader)

	for {
		line, err := reader.ReadString('\n') // Read until newline
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "data: [DONE]" {
			return nil, io.EOF // End of stream
		}
		if len(line) > 6 && line[:6] == "data: " {
			trimmed := line[6:] // Trim the "data: " prefix
			var response StreamChatCompletionResponse
			if err := json.Unmarshal([]byte(trimmed), &response); err != nil {
				return nil, fmt.Errorf("unmarshal error: %w, raw data: %s", err, trimmed)
			}
			if response.Usage == nil {
				response.Usage = &StreamUsage{}
			}
			return &response, nil
		}
	}
}

// Close terminates the stream.
func (s *chatCompletionStream) Close() error {
	s.cancel()
	return s.resp.Body.Close()
}
