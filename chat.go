package deepseek

import (
	"errors"
)

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

var (
	ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method")
	ErrUnexpectedResponseFormat         = errors.New("unexpected response format")
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
}
