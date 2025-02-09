package deepseek_examples

import (
	"fmt"
	"log"
	"time"

	"github.com/cohesion-org/deepseek-go"
)

func NewClientWithOptions() {
	// Initialize client with a custom BaseURL and Timeout
	client1, err := deepseek.NewClientWithOptions("your-api-key",
		deepseek.WithBaseURL("https://custom-api.com/"),
		deepseek.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatalf("Error creating client1: %v", err)
	}
	fmt.Printf("Client1 initialized with BaseURL: %s and Timeout: %v\n", client1.BaseURL, client1.Timeout)

	// Initialize client with only BaseURL
	client2, err := deepseek.NewClientWithOptions("your-api-key",
		deepseek.WithBaseURL("https://alternative-api.com/"),
	)
	if err != nil {
		log.Fatalf("Error creating client2: %v", err)
	}
	fmt.Printf("Client2 initialized with BaseURL: %s and Timeout: %v\n", client2.BaseURL, client2.Timeout)

	// Initialize client with only Timeout
	client3, err := deepseek.NewClientWithOptions("your-api-key",
		deepseek.WithTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("Error creating client3: %v", err)
	}
	fmt.Printf("Client3 initialized with BaseURL: %s and Timeout: %v\n", client3.BaseURL, client3.Timeout)

	// Initialize client with Timeout set via string parsing
	client4, err := deepseek.NewClientWithOptions("your-api-key",
		deepseek.WithTimeoutString("45s"),
	)
	if err != nil {
		log.Fatalf("Error creating client4: %v", err)
	}
	fmt.Printf("Client4 initialized with BaseURL: %s and Timeout: %v\n", client4.BaseURL, client4.Timeout)

	// Initialize client with default settings
	client5, err := deepseek.NewClientWithOptions("your-api-key")
	if err != nil {
		log.Fatalf("Error creating client5: %v", err)
	}
	fmt.Printf("Client5 initialized with BaseURL: %s and Timeout: %v\n", client5.BaseURL, client5.Timeout)
}
