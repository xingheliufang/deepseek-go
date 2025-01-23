# Deepseek-Go

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Deepseek-Go is a Go-based API wrapper for the [Deepseek](https://deepseek.com) platform. It provides a clean and type-safe interface to interact with Deepseek's AI features, including chat completions with streaming, token usage tracking, and more.

This library is designed for developers building Go applications that require seamless integration with Deepseek AI.

## Features

- **Chat Completion**: Easily send chat messages and receive responses from Deepseek's AI models. It also supports streaming.
- **Modular Design**: The library is structured into reusable components for building, sending, and handling requests and responses.
- **MIT License**: Open-source and free for both personal and commercial use.

## Installation

To use Deepseek-Go, ensure you have Go installed, and run:

```sh
go get github.com/cohesion-org/deepseek-go
```

## Getting Started

Here's a quick example of how to use the library:

### Prerequisites

Before using the library, ensure you have:
- A valid Deepseek API key.
- Go installed on your system.


<details open>
<summary> Chat </summary>

### Example for chatting with deepseek

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
	constants "github.com/cohesion-org/deepseek-go/constants"
)

func main() {
	// Set up the Deepseek client
    client := deepseek.NewClient(os.Getenv("DEEPSEEK_KEY"))

	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleSystem, Content: "Answer every question using slang."},
			{Role: constants.ChatMessageRoleUser, Content: "Which is the tallest mountain in the world?"},
		},
	}

	// Send the request and handle the response
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print the response
	fmt.Println("Response:", response.Choices[0].Message.Content)
}
```

Save this code to a file (e.g., `main.go`), and run it:
```sh
go run main.go
```
</details>

<details >
	<summary> Sedning other params like Temp, Stop </summary>
	<strong> You just need to extend the ChatCompletionMessage with the supported parameters. </strong>

```go
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "What is the meaning of deepseek"},
			{Role: constants.ChatMessageRoleSystem, Content: "Answer every question using slang"},
		},
		Temperature: 1.0,
		Stop:        []string{"yo", "hello"},
		ResponseFormat: &deepseek.ResponseFormat{
			Type: "text",
		},
	}
```

</details>

<details>
<summary> Chat with Streaming </summary>

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
	constants "github.com/cohesion-org/deepseek-go/constants"
)

func main() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_KEY"))
	request := &deepseek.StreamChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: constants.ChatMessageRoleUser, Content: "Just testing if the streaming feature is working or not!"},
		},
		Stream: true,
	}
	ctx := context.Background()

	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		log.Fatalf("ChatCompletionStream error: %v", err)
	}
	var fullMessage string
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			break
		}
		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break
		}
		for _, choice := range response.Choices {
			fullMessage += choice.Delta.Content // Accumulate chunk content
			log.Println(choice.Delta.Content)
		}
	}
	log.Println("The full message is: ", fullMessage)
}
```
</details>

<details>
<summary> Get the balance(s) of the user. </summary>

```go
package main

import (
	"context"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

func main() {
	client := deepseek.NewClient(os.Getenv("DEEPSEEK_KEY"))
	ctx := context.Background()
	balance, err := deepseek.GetBalance(client, ctx)
	if err != nil {
		log.Fatalf("Error getting balance: %v", err)
	}

	if balance == nil {
		log.Fatalf("Balance is nil")
	}

	if len(balance.BalanceInfos) == 0 {
		log.Fatalf("No balance information returned")
	}
	log.Printf("%+v\n", balance)
}
```
</details>


---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Credits
- **`chat.go` Inspiration**: Adapted from [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai/tree/master).

---

Feel free to contribute, open issues, or submit PRs to help improve Deepseek-Go! Let us know if you encounter any issues.