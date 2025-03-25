package deepseek_test

import (
	"encoding/json"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	"github.com/stretchr/testify/assert"
)

func TestEstimateTokenCount(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		minCount int // Since these are estimates, we test for minimum expected tokens
	}{
		{
			name:     "english text",
			text:     "Hello, world!",
			minCount: 3, // At least 3 tokens for "Hello", ",", "world!"
		},
		{
			name:     "chinese text",
			text:     "你好世界",
			minCount: 2, // At least 2 tokens for 4 Chinese characters
		},
		{
			name:     "mixed text",
			text:     "Hello 世界!",
			minCount: 2, // At least 2 tokens
		},
		{
			name:     "empty text",
			text:     "",
			minCount: 1, // Minimum 1 token
		},
		{
			name:     "numbers and symbols",
			text:     "123 !@#",
			minCount: 2, // At least 2 tokens
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate := deepseek.EstimateTokenCount(tt.text)
			assert.NotNil(t, estimate)
			assert.GreaterOrEqual(t, estimate.EstimatedTokens, tt.minCount)
		})
	}
}

func TestEstimateTokensFromMessages(t *testing.T) {
	tests := []struct {
		name     string
		request  *deepseek.ChatCompletionRequest
		minCount int // Since these are estimates, we test for minimum expected tokens
	}{
		{
			name: "single message",
			request: &deepseek.ChatCompletionRequest{
				Messages: []deepseek.ChatCompletionMessage{
					{
						Role:    constants.ChatMessageRoleUser,
						Content: "Hello!",
					},
				},
			},
			minCount: 4, // At least 4 tokens (2 for role + content tokens)
		},
		{
			name: "multiple messages",
			request: &deepseek.ChatCompletionRequest{
				Messages: []deepseek.ChatCompletionMessage{
					{
						Role:    constants.ChatMessageRoleSystem,
						Content: "You are a helpful assistant.",
					},
					{
						Role:    constants.ChatMessageRoleUser,
						Content: "Hi!",
					},
				},
			},
			minCount: 8, // At least 8 tokens (2 per role + content tokens)
		},
		{
			name: "with function definition",
			request: &deepseek.ChatCompletionRequest{
				Messages: []deepseek.ChatCompletionMessage{
					{
						Role:    constants.ChatMessageRoleUser,
						Content: "Get weather",
					},
				},
				Tools: []deepseek.Tool{
					{
						Function: deepseek.Function{
							Name:        "get_weather",
							Description: "Get weather information",
							Parameters: json.RawMessage(`
								{
									"type": "object",
									"properties": {
										"location": {
											"type": "string",
											"description": "The city and state, e.g. San Francisco, CA"
										}
									},
									"required": [
										"location"
									]
								}`,
							),
						},
					},
				},
			},
			minCount: 12, // At least 12 tokens (message + function definition tokens)
		},
		{
			name: "empty request",
			request: &deepseek.ChatCompletionRequest{
				Messages: []deepseek.ChatCompletionMessage{},
			},
			minCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate := deepseek.EstimateTokensFromMessages(tt.request)
			assert.NotNil(t, estimate)
			assert.GreaterOrEqual(t, estimate.EstimatedTokens, tt.minCount)
		})
	}
}
