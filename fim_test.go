package deepseek_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFIMCompletion(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	tests := []struct {
		name        string
		req         *deepseek.FIMCompletionRequest
		wantErr     bool
		validateRes func(t *testing.T, res *deepseek.FIMCompletionResponse)
	}{
		{
			name: "basic FIM completion",
			req: &deepseek.FIMCompletionRequest{
				Model:  deepseek.DeepSeekChat,
				Prompt: "func main() {\n    fmt.Println(\"hel",
			},
			wantErr: false,
			validateRes: func(t *testing.T, res *deepseek.FIMCompletionResponse) {
				assert.NotEmpty(t, res.Choices[0].Text)
			},
		},
		{
			name: "empty prompt",
			req: &deepseek.FIMCompletionRequest{
				Model:  deepseek.DeepSeekChat,
				Prompt: "",
			},
			wantErr: true,
		},
		{
			name: "invalid model",
			req: &deepseek.FIMCompletionRequest{
				Model:  "invalid-model-123",
				Prompt: "some code",
			},
			wantErr: true,
		},
		{
			name: "max tokens exceeded",
			req: &deepseek.FIMCompletionRequest{
				Model:     deepseek.DeepSeekChat,
				Prompt:    "long prompt " + strings.Repeat("test ", 1000),
				MaxTokens: 5000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			resp, err := client.CreateFIMCompletion(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp)

			if tt.req != nil {
				// Validate common response structure
				assert.NotEmpty(t, resp.ID)
				assert.NotEmpty(t, resp.Created)
				assert.Equal(t, "text_completion", resp.Object)
				assert.Equal(t, tt.req.Model, resp.Model)
				assert.NotEmpty(t, resp.Choices)
				assert.NotNil(t, resp.Usage)
			}
			// Validate specific test case expectations
			if tt.validateRes != nil {
				tt.validateRes(t, resp)
			}
		})
	}
}

func TestCreateFIMCompletionWithParameters(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	tests := []struct {
		name        string
		req         *deepseek.FIMCompletionRequest
		wantErr     bool
		validateRes func(t *testing.T, res *deepseek.FIMCompletionResponse)
	}{
		{
			name: "FIM completion with temperature and top_p",
			req: &deepseek.FIMCompletionRequest{
				Model:       deepseek.DeepSeekChat,
				Prompt:      "func main() {\n    fmt.Println(\"hel",
				Temperature: 0.5,
				TopP:        0.9,
			},
			wantErr: false,
			validateRes: func(t *testing.T, res *deepseek.FIMCompletionResponse) {
				assert.NotEmpty(t, res.Choices[0].Text)
			},
		},
		{
			name: "FIM completion with invalid temperature",
			req: &deepseek.FIMCompletionRequest{
				Model:       deepseek.DeepSeekChat,
				Prompt:      "func main() {\n    fmt.Println(\"hel",
				Temperature: 2.5, // Invalid temperature
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			resp, err := client.CreateFIMCompletion(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp)

			if tt.req != nil {
				// Validate common response structure (same as in TestCreateFIMCompletion)
				assert.NotEmpty(t, resp.ID)
				assert.NotEmpty(t, resp.Created)
				assert.Equal(t, "text_completion", resp.Object)
				assert.Equal(t, tt.req.Model, resp.Model)
				assert.NotEmpty(t, resp.Choices)
				assert.NotNil(t, resp.Usage)
			}
			// Validate specific test case expectations
			if tt.validateRes != nil {
				tt.validateRes(t, resp)
			}
		})
	}
}

func TestFIMCompletionResponseStructure(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	req := &deepseek.FIMCompletionRequest{
		Model:  deepseek.DeepSeekChat,
		Prompt: "func main() {\n    fmt.Println(\"hel",
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	resp, err := client.CreateFIMCompletion(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Check overall response structure
	assert.NotEmpty(t, resp.ID)
	assert.NotEmpty(t, resp.Object)
	assert.NotEmpty(t, resp.Created)
	assert.Equal(t, "text_completion", resp.Object)
	assert.Equal(t, req.Model, resp.Model)
	assert.NotEmpty(t, resp.Choices)
	assert.NotNil(t, resp.Usage)

	// Check choices array
	for _, choice := range resp.Choices {
		assert.NotEmpty(t, choice.Text)
		assert.NotNil(t, choice.Index) // Or assert.GreaterOrEqual(t, choice.Index, 0)
		// LogProbs can be nil, so just check that it's present or handle it if you expect values.
		assert.NotEmpty(t, choice.FinishReason) // Check for valid values like "stop", "length"

		//Optional, if logprobs are requested
		if req.LogProbs > 0 {
			assert.NotNil(t, choice.LogProbs)
		}
	}

	// Check usage structure
	assert.GreaterOrEqual(t, resp.Usage.PromptTokens, 0)
	assert.GreaterOrEqual(t, resp.Usage.CompletionTokens, 0)
	assert.GreaterOrEqual(t, resp.Usage.TotalTokens, 0)
	assert.Equal(t, resp.Usage.PromptTokens+resp.Usage.CompletionTokens, resp.Usage.TotalTokens) // Check if the total is correct

}

func TestFIMCompletionResponseWithLogProbs(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	req := &deepseek.FIMCompletionRequest{
		Model:    deepseek.DeepSeekChat,
		Prompt:   "func main() {\n    fmt.Println(\"hel",
		LogProbs: 1, // Request log probabilities
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	resp, err := client.CreateFIMCompletion(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Check choices array
	for _, choice := range resp.Choices {
		assert.NotEmpty(t, choice.Text)
		assert.NotNil(t, choice.Index)
		assert.NotNil(t, choice.LogProbs) // Now LogProbs should definitely be present

		logProbsMap, ok := choice.LogProbs.(map[string]interface{})
		if ok {
			assert.NotEmpty(t, logProbsMap) // Check that the map is not empty
		} else {
			t.Errorf("Unexpected type for LogProbs: %T", choice.LogProbs)
		}

		assert.NotEmpty(t, choice.FinishReason)
	}
}
