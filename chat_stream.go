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

// StreamChatCompletionMessage represents a single message in a chat completion stream.
type StreamChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionStream is an interface for receiving streaming chat completion responses.
type ChatCompletionStream interface {
	Recv() (*StreamChatCompletionResponse, error)
	Close() error
}

// chatCompletionStream implements the ChatCompletionStream interface.
type chatCompletionStream struct {
	ctx    context.Context    // Context for cancellation.
	cancel context.CancelFunc // Cancel function for the context.
	resp   *http.Response     // HTTP response from the API call.
	reader *bufio.Reader      // Reader for the response body.
}

// StreamDelta represents a delta in the chat completion stream.
type StreamDelta struct {
	Role             string     `json:"role,omitempty"`              // Role of the message.
	Content          string     `json:"content"`                     // Content of the message.
	ReasoningContent string     `json:"reasoning_content,omitempty"` // Reasoning content of the message.
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`        // Optional tool calls related to the message.
}

// StreamChoices represents a choice in the chat completion stream.
type StreamChoices struct {
	Index        int         `json:"index"`                   // Index of the choice.
	Delta        StreamDelta `json:"delta"`                   // Delta information for the choice.
	FinishReason string      `json:"finish_reason,omitempty"` // Reason for finishing the generation.
	Logprobs     Logprobs    `json:"logprobs,omitempty"`      // Log probabilities for the generated tokens.
}

// StreamChatCompletionResponse represents a single response from a streaming chat completion API call.
type StreamChatCompletionResponse struct {
	ID      string          `json:"id"`              // ID of the response.
	Object  string          `json:"object"`          // Type of object.
	Created int64           `json:"created"`         // Creation timestamp.
	Model   string          `json:"model"`           // Model used.
	Choices []StreamChoices `json:"choices"`         // Choices generated.
	Usage   *Usage          `json:"usage,omitempty"` // Usage statistics (optional).
}

// Recv receives the next response from the stream.
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
			return &response, nil
		}
	}
}

// Close terminates the stream.
func (s *chatCompletionStream) Close() error {
	s.cancel()
	err := s.resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}
	return nil
}
