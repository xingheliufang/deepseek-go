package deepseek

import (
	"encoding/json"
	"errors"

	"github.com/cohesion-org/deepseek-go/constants"
)

const (
	// ChatMessageRoleSystem is the role of a system message
	ChatMessageRoleSystem = constants.ChatMessageRoleSystem
	// ChatMessageRoleUser is the role of a user message
	ChatMessageRoleUser = constants.ChatMessageRoleUser
	// ChatMessageRoleAssistant is the role of an assistant message
	ChatMessageRoleAssistant = constants.ChatMessageRoleAssistant
	// ChatMessageRoleTool is the role of a tool message
	ChatMessageRoleTool = constants.ChatMessageRoleTool
)

var (
	// ErrChatCompletionStreamNotSupported is returned when streaming is not supported with the method.
	ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method")
	// ErrChatCompletionRequestNil is returned when the request is nil.
	ErrUnexpectedResponseFormat = errors.New("unexpected response format")
)

// ChatCompletionMessage represents a single message in a chat completion conversation.
type ChatCompletionMessage struct {
	Role             string     `json:"role"`                        // The role of the message sender, e.g., "user", "assistant", "system".
	Content          string     `json:"content"`                     // The content of the message.
	Prefix           bool       `json:"prefix,omitempty"`            // The prefix of the message (optional) for Chat Prefix Completion [Beta Feature].
	ReasoningContent string     `json:"reasoning_content,omitempty"` // The reasoning content of the message (optional) when using the reasoner model with Chat Prefix Completion. When using this feature, the Prefix parameter must be set to true.
	ToolCallID       string     `json:"tool_call_id,omitempty"`      // Tool call that this message is responding to.
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`        // Optional tool calls.
}

// Function defines the structure of a function tool.
type Function struct {
	Name        string          `json:"name"`                 // The name of the function (required).
	Description string          `json:"description"`          // A description of the function (required).
	Parameters  json.RawMessage `json:"parameters,omitempty"` // The parameters of the function (optional).
}

// Tool defines the structure for a tool.
type Tool struct {
	Type     string   `json:"type"`     // The type of the tool, e.g., "function" (required).
	Function Function `json:"function"` // The function details (required).
}

// ResponseFormat defines the structure for the response format.
type ResponseFormat struct {
	Type string `json:"type"` // The desired response format, either "text" or "json_object".
}

type Stop []string

func (s *Stop) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var xs any
	if err := json.Unmarshal(data, &xs); err != nil {
		return err
	}
	switch stop := xs.(type) {
	case string:
		*s = []string{stop}
	case []interface{}:
		var stops []string
		for _, x := range stop {
			if v, ok := x.(string); ok {
				stops = append(stops, v)
			} else {
				return errors.New("invalid type for stop: expected string or array of strings")
			}
		}
		*s = stops
	default:
		return errors.New("invalid type for stop: expected string or array of strings")
	}
	return nil
}

func (s *Stop) MarshalJSON() ([]byte, error) {
	if len(*s) == 0 {
		return json.Marshal(nil)
	}
	if len(*s) == 1 {
		return json.Marshal((*s)[0])
	}
	return json.Marshal([]string(*s))
}

// StreamOptions provides options for streaming chat completion responses.
type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"` // Whether to include usage information in the stream. The API returns the usage sometimes even if this is set to false.
}

// ChatCompletionRequest defines the structure for a chat completion request.
type ChatCompletionRequest struct {
	Model            string                  `json:"model"`                       // The ID of the model to use (required).
	Messages         []ChatCompletionMessage `json:"messages"`                    // A list of messages comprising the conversation (required).
	FrequencyPenalty *float32                `json:"frequency_penalty,omitempty"` // Penalty for new tokens based on their frequency in the text so far (optional).
	MaxTokens        *int                    `json:"max_tokens,omitempty"`        // The maximum number of tokens to generate in the chat completion (optional).
	PresencePenalty  *float32                `json:"presence_penalty,omitempty"`  // Penalty for new tokens based on their presence in the text so far (optional).
	Temperature      *float32                `json:"temperature,omitempty"`       // The sampling temperature, between 0 and 2 (optional).
	TopP             *float32                `json:"top_p,omitempty"`             // The nucleus sampling parameter, between 0 and 1 (optional).
	ResponseFormat   *ResponseFormat         `json:"response_format,omitempty"`   // The desired response format (optional).
	Stop             Stop                    `json:"stop,omitempty"`              // A list of sequences where the model should stop generating further tokens (optional).
	Stream           *bool                   `json:"stream,omitempty"`            //Comments: Defaults to true, since it's "STREAM"
	StreamOptions    *StreamOptions          `json:"stream_options,omitempty"`    // Optional: Stream options for the request.
	Tools            []Tool                  `json:"tools,omitempty"`             // A list of tools the model may use (optional).
	ToolChoice       *OneOfToolChoice        `json:"tool_choice,omitempty"`       // Controls which (if any) tool is called by the model (optional).
	LogProbs         *bool                   `json:"logprobs,omitempty"`          // Whether to return log probabilities of the most likely tokens (optional).
	TopLogProbs      *int                    `json:"top_logprobs,omitempty"`      // The number of top most likely tokens to return log probabilities for (optional).
	JSONMode         *bool                   `json:"json,omitempty"`              // [deepseek-go feature] Optional: Enable JSON mode. If you're using the JSON mode, please mention "json" anywhere in your prompt, and also include the JSON schema in the request.
}

type ToolChoice interface {
	IsToolChoice()
}

type ChatCompletionToolChoice string

func (ChatCompletionToolChoice) IsToolChoice() {}

// ToolChoiceFunction defines the function details within ToolChoice.
type ToolChoiceFunction struct {
	Name string `json:"name"` // The name of the function to call (required).
}

type ChatCompletionNamedToolChoice struct {
	Type     string             `json:"type"`
	Function ToolChoiceFunction `json:"function"`
}

func (ChatCompletionNamedToolChoice) IsToolChoice() {}

type OneOfToolChoice struct {
	ToolChoice
}

func (o *OneOfToolChoice) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var x any
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	switch toolChoice := x.(type) {
	case string:
		o.ToolChoice = ChatCompletionToolChoice(toolChoice)
	default:
		var namedToolChoice ChatCompletionNamedToolChoice
		if err := json.Unmarshal(data, &namedToolChoice); err != nil {
			return err
		}
		o.ToolChoice = namedToolChoice
	}

	return nil
}

func (o *OneOfToolChoice) MarshalJSON() ([]byte, error) {
	if o.ToolChoice == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(o.ToolChoice)
}
