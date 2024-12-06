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
	Content string `json:"content"`
}

type StreamChatCompletionRequest struct {
	Stream           bool                    `json:"stream,omitempty"`            //Comments: Defaults to true, since it's "STREAM"
	Model            string                  `json:"model"`                       // Required: Model ID, e.g., "deepseek-chat"
	Messages         []ChatCompletionMessage `json:"messages"`                    // Required: List of messages
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"` // Optional: Frequency penalty, >= -2 and <= 2
	MaxTokens        int                     `json:"max_tokens,omitempty"`        // Optional: Maximum tokens, > 1
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`  // Optional: Presence penalty, >= -2 and <= 2
	Temperature      float32                 `json:"temperature,omitempty"`       // Optional: Sampling temperature, <= 2
	TopP             float32                 `json:"top_p,omitempty"`             // Optional: Nucleus sampling parameter, <= 1
	ResponseFormat   *ResponseFormat         `json:"response_format,omitempty"`   // Optional: Custom response format: just don't try, it breaks rn ;)
	Stop             []string                `json:"stop,omitempty"`              // Optional: Stop signals
	Tools            []Tools                 `json:"tools,omitempty"`             // Optional: List of tools
	LogProbs         bool                    `json:"logprobs,omitempty"`          // Optional: Enable log probabilities
	TopLogProbs      int                     `json:"top_logprobs,omitempty"`      // Optional: Number of top tokens with log probabilities, <= 20
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
	reader := s.reader
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
