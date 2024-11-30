# Deepseek-Go

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Deepseek-Go is a Go-based API wrapper for the [Deepseek AI](https://deepseek.ai) platform. It provides a clean and type-safe interface to interact with Deepseek's AI features, including chat completions, token usage tracking, and more.

This library is designed for developers building Go applications that require seamless integration with Deepseek AI.

## Features

- **Chat Completion**: Easily send chat messages and receive responses from Deepseek's AI models.
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

### Example Code

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	deepseek "github.com/cohesion-org/deepseek-go"
)

func main() {
	// Set up the Deepseek client
	client := deepseek.NewClient("your deepseek api key")

	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekCoder,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleSystem, Content: "Answer every question using slang."},
			{Role: deepseek.ChatMessageRoleUser, Content: "Which is the tallest mountain in the world?"},
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

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## TODO

- [ ] Make streaming possible.
- [ ] Add examples for token usage tracking.
- [ ] Improve error handling in the wrapper.
- [ ] Add support for other Deepseek endpoints.
- [ ] Write comprehensive tests.

---

## Credits
- **`chat.go` Inspiration**: Adapted from [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai/tree/master).

---

Feel free to contribute, open issues, or submit PRs to help improve Deepseek-Go! Let us know if you encounter any issues.