package Utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type Choice struct {
	Index        int       `json:"index"`
	Message      Message   `json:"message"`
	LogProbs     *LogProbs `json:"logprobs,omitempty"`
	FinishReason string    `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LogProbs struct {
	Tokens        []string             `json:"tokens,omitempty"`
	TokenLogProbs []float64            `json:"token_logprobs,omitempty"`
	TopLogProbs   []map[string]float64 `json:"top_logprobs,omitempty"`
}

type Usage struct {
	PromptTokens          int `json:"prompt_tokens"`
	CompletionTokens      int `json:"completion_tokens"`
	TotalTokens           int `json:"total_tokens"`
	PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
}

func HandleResponse(resp *http.Response) (*ChatCompletionResponse, error) {
	body, err := io.ReadAll(resp.Body) //never re read your body, do it only once. io clears when reading

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	var parsedResponse ChatCompletionResponse
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return &parsedResponse, nil
}
