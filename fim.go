package deepseek

// FIMCompletionRequest represents the request body for a Fill-In-the-Middle (FIM) completion.
type FIMCompletionRequest struct {
	Model            string   `json:"model"`                       // Model name to use for completion.
	Prompt           string   `json:"prompt"`                      // The prompt to start the completion from.
	Suffix           string   `json:"suffix,omitempty"`            // Optional: The suffix to complete the prompt with.
	MaxTokens        int      `json:"max_tokens,omitempty"`        // Optional: Maximum tokens to generate, > 1 and <= 4000.
	Temperature      float64  `json:"temperature,omitempty"`       // Optional: Sampling temperature, between 0 and 1.
	TopP             float64  `json:"top_p,omitempty"`             // Optional: Nucleus sampling probability threshold.
	N                int      `json:"n,omitempty"`                 // Optional: Number of completions to generate.
	LogProbs         int      `json:"logprobs,omitempty"`          // Optional: Number of log probabilities to return.
	Echo             bool     `json:"echo,omitempty"`              // Optional: Whether to echo the prompt in the completion.
	Stop             []string `json:"stop,omitempty"`              // Optional: List of stop sequences.
	PresencePenalty  float64  `json:"presence_penalty,omitempty"`  // Optional: Penalty for new tokens based on their presence in the text so far.
	FrequencyPenalty float64  `json:"frequency_penalty,omitempty"` // Optional: Penalty for new tokens based on their frequency in the text so far.
}

// FIMCompletionResponse represents the response body for a Fill-In-the-Middle (FIM) completion.
type FIMCompletionResponse struct {
	ID      string `json:"id"`      // Unique ID for the completion.
	Object  string `json:"object"`  // The object type, e.g., "text_completion".
	Created int    `json:"created"` // Timestamp of when the completion was created.
	Model   string `json:"model"`   // Model used for the completion.
	Choices []struct {
		Text         string `json:"text"`          // The generated completion text.
		Index        int    `json:"index"`         // Index of the choice.
		LogProbs     any    `json:"logprobs"`      // Log probabilities of the generated tokens (if requested).  Type 'any' because it can be null.
		FinishReason string `json:"finish_reason"` // Reason for finishing the completion, e.g., "stop", "length".
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`     // Number of tokens in the prompt.
		CompletionTokens int `json:"completion_tokens"` // Number of tokens in the completion.
		TotalTokens      int `json:"total_tokens"`      // Total number of tokens used.
	} `json:"usage"`
}

// FIMStreamCompletionRequest represents the request body for a streaming Fill-In-the-Middle (FIM) completion.
// It's similar to FIMCompletionRequest but includes a `Stream` field.
type FIMStreamCompletionRequest struct {
	Model            string        `json:"model"`                       // Model name to use for completion.
	Prompt           string        `json:"prompt"`                      // The prompt to start the completion from.
	Stream           bool          `json:"stream"`                      // Whether to stream the completion.  This is the key difference.
	StreamOptions    StreamOptions `json:"stream_options,omitempty"`    // Optional: Options for streaming the completion.
	Suffix           string        `json:"suffix,omitempty"`            // Optional: The suffix to complete the prompt with.
	MaxTokens        int           `json:"max_tokens,omitempty"`        // Optional: Maximum tokens to generate, > 1 and <= 4000.
	Temperature      float64       `json:"temperature,omitempty"`       // Optional: Sampling temperature, between 0 and 1.
	TopP             float64       `json:"top_p,omitempty"`             // Optional: Nucleus sampling probability threshold.
	N                int           `json:"n,omitempty"`                 // Optional: Number of completions to generate.
	LogProbs         int           `json:"logprobs,omitempty"`          // Optional: Number of log probabilities to return.
	Echo             bool          `json:"echo,omitempty"`              // Optional: Whether to echo the prompt in the completion.
	Stop             []string      `json:"stop,omitempty"`              // Optional: List of stop sequences.
	PresencePenalty  float64       `json:"presence_penalty,omitempty"`  // Optional: Penalty for new tokens based on their presence in the text so far.
	FrequencyPenalty float64       `json:"frequency_penalty,omitempty"` // Optional: Penalty for new tokens based on their frequency in the text so far.
}

// // ChatCompletionStream is an interface for receiving streaming chat completion responses.
// type FIMChatCompletionStream interface {
// 	FIMRecv() (*StreamChatCompletionResponse, error)
// 	Close() error
// }

// // chatCompletionStream implements the ChatCompletionStream interface.
// type fimchatCompletionStream struct {
// 	ctx    context.Context    // Context for cancellation.
// 	cancel context.CancelFunc // Cancel function for the context.
// 	resp   *http.Response     // HTTP response from the API call.
// 	reader *bufio.Reader      // Reader for the response body.
// }

// // FIMRecv receives the next response from the FIMstream.
// func (s *fimchatCompletionStream) FIMRecv() (*StreamChatCompletionResponse, error) {
// 	reader := s.reader
// 	for {
// 		line, err := reader.ReadString('\n') // Read until newline
// 		if err != nil {
// 			if err == io.EOF {
// 				return nil, io.EOF
// 			}
// 			return nil, fmt.Errorf("error reading stream: %w", err)
// 		}

// 		line = strings.TrimSpace(line)
// 		if line == "data: [DONE]" {
// 			return nil, io.EOF // End of stream
// 		}
// 		if len(line) > 6 && line[:6] == "data: " {
// 			trimmed := line[6:] // Trim the "data: " prefix
// 			var response StreamChatCompletionResponse
// 			if err := json.Unmarshal([]byte(trimmed), &response); err != nil {
// 				return nil, fmt.Errorf("unmarshal error: %w, raw data: %s", err, trimmed)
// 			}
// 			if response.Usage == nil {
// 				response.Usage = &StreamUsage{}
// 			}
// 			return &response, nil
// 		}
// 	}
// }
