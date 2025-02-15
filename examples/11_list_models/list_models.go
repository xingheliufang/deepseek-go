package deepseek_examples

import (
	"context"
	"fmt"
	"log"

	deepseek "github.com/cohesion-org/deepseek-go"
)

// ListAllModels demonstrates how to list all models.
func ListAllModels() {
	client := deepseek.NewClient("DEEPSEEK_API_KEY")
	ctx := context.Background()
	models, err := deepseek.ListAllModels(client, ctx)
	if err != nil {
		log.Fatalf("Error listing models: %v", err)
	}
	fmt.Printf("\n%+v\n", models)
}
