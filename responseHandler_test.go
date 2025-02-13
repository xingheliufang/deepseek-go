package deepseek_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleChatCompletionResponse(t *testing.T) {
	validResponse := deepseek.ChatCompletionResponse{
		ID:      "chat-123",
		Object:  "chat.completion",
		Created: 1677858242,
		Model:   "deepseek-chat",
		Choices: []deepseek.Choice{
			{
				Index: 0,
				Message: deepseek.Message{
					Role:    "assistant",
					Content: "Hello! How can I help you today?",
				},
				FinishReason: "stop",
			},
		},
		Usage: deepseek.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}

	tests := []struct {
		name        string
		response    *http.Response
		want        *deepseek.ChatCompletionResponse
		wantErr     string
		errContains string
	}{
		{
			name: "valid response",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewReader([]byte(`{
					"id": "chat-123",
					"object": "chat.completion",
					"created": 1677858242,
					"model": "deepseek-chat",
					"choices": [{
						"index": 0,
						"message": {
							"role": "assistant",
							"content": "Hello! How can I help you today?"
						},
						"finish_reason": "stop"
					}],
					"usage": {
						"prompt_tokens": 10,
						"completion_tokens": 20,
						"total_tokens": 30
					}
				}`))),
			},
			want: &validResponse,
		},
		{
			name: "invalid JSON",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"invalid": `))),
			},
			wantErr:     "failed to parse response JSON",
			errContains: "unexpected end of JSON input",
		},
		{
			name: "empty body",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte{})),
			},
			wantErr:     "failed to parse response JSON",
			errContains: "empty response body",
		},
		{
			name: "read error",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(&errorReader{
					err: errors.New("read error"),
				}),
			},
			wantErr:     "failed to read response body",
			errContains: "read error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.response.Body.Close()

			resp, err := deepseek.HandleChatCompletionResponse(tt.response)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, resp)
		})
	}
}

type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (int, error) {
	return 0, r.err
}

func TestResponseStructureValidation(t *testing.T) {
	t.Run("missing required fields", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewReader([]byte(`{
				"created": 1677858242,
				"choices": [{}]
			}`))),
		}
		defer resp.Body.Close()

		_, err := deepseek.HandleChatCompletionResponse(resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid response")
	})

	t.Run("unexpected field types", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewReader([]byte(`{
				"id": 123,
				"object": "chat.completion",
				"created": "invalid",
				"model": "deepseek-chat",
				"choices": [{"index": "zero"}],
				"usage": {"prompt_tokens": "ten"}
			}`))),
		}
		defer resp.Body.Close()

		_, err := deepseek.HandleChatCompletionResponse(resp)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse response JSON")
	})
}

func TestErrorWrapping(t *testing.T) {
	t.Run("read error wrapping", func(t *testing.T) {
		resp := &http.Response{
			Body: io.NopCloser(&errorReader{err: errors.New("custom error")}),
		}
		_, err := deepseek.HandleChatCompletionResponse(resp)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read response body")
	})

	t.Run("json error wrapping", func(t *testing.T) {
		resp := &http.Response{
			Body: io.NopCloser(bytes.NewReader([]byte(`{"invalid": `))),
		}
		_, err := deepseek.HandleChatCompletionResponse(resp)
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse response JSON")
	})
}
