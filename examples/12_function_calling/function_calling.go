package deepseek_examples

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cohesion-org/deepseek-go"
)

// When you provide multiple utility functions, the model may call multiple utility functions
// in one request and multiple times in one conversation. Specifically, the model may pass the
// result of the previously called utility function as input to the second called utility function.
//
// Example:
//   User:      What the weather and relative humidity are at the current location?
//   Assistant: call GetLocation()
//   Tool:      "Beijing"
//   Assistant: call GetTemperature("Beijing") and call GetRelativeHumidity("Beijing")
//   Tool:      "11℃", "35%"
//   Assistant: The current temperature and humidity in Beijing are 11℃ and 35%.
//
// Yes, this happens in one request, so you can use it to solve complex problems.

var toolGetTime = deepseek.Tool{
	Type: "function",
	Function: deepseek.Function{
		Name: "GetTime",
		Description: "" +
			"Get the current date and time. The returned time string format is RFC3339. " +
			"Be careful not to abuse this function unless you really need to get the real world time.",
	},
}

func onGetTime() string {
	s := time.Now().Format(time.RFC3339)
	return "current time: " + s
}

func FunctionCalling() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_API_KEY"))

	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: "What time is it now?"},
		},
		Tools: []deepseek.Tool{toolGetTime},
	}
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// will be empty
	fmt.Println("response:", response.Choices[0].Message.Content)

	// one tool call request
	fmt.Println("tool calls:", response.Choices[0].Message.ToolCalls)

	msg := response.Choices[0].Message
	toolCalls := msg.ToolCalls

	question := deepseek.ChatCompletionMessage{
		Role:      deepseek.ChatMessageRoleAssistant,
		Content:   msg.Content,
		ToolCalls: toolCalls,
	}
	answer := deepseek.ChatCompletionMessage{
		Role:       deepseek.ChatMessageRoleTool,
		Content:    onGetTime(),
		ToolCallID: toolCalls[0].ID,
	}

	messages := request.Messages
	messages = append(messages, question, answer)
	toolReq := &deepseek.ChatCompletionRequest{
		Model:    request.Model,
		Messages: messages,

		// It is not recommended to use it unless it is a special case.
		// The official said that they are actively fixing the problem
		// of infinite loop calls.

		// Not using this field will force the model to call the
		// utility function only once per conversation.

		// Don't try to delete only the utility functions which may cause
		// infinite loop calls. I have tried this and the model will still
		// call the deleted utility functions unless you delete all the
		// utility functions, like the commented out code.

		// Tools: request.Tools, // This is the key to implement chain calls.
	}

	response, err = client.CreateChatCompletion(ctx, toolReq)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// will return the current time
	fmt.Println("response:", response.Choices[0].Message.Content)
}
