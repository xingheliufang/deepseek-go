package deepseek_test

import (
	"context"
	"testing"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAllModels(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	resp, err := deepseek.ListAllModels(client, ctx)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Verify response structure
	assert.Equal(t, "list", resp.Object)
	assert.NotEmpty(t, resp.Data)

	// Verify model details
	for _, model := range resp.Data {
		assert.NotEmpty(t, model.ID)
		assert.Equal(t, "model", model.Object)
		assert.Equal(t, "deepseek", model.OwnedBy)

		// Verify known models exist in constants.go
		if model.ID == deepseek.DeepSeekChat ||
			model.ID == deepseek.DeepSeekCoder ||
			model.ID == deepseek.DeepSeekReasoner {
			assert.Contains(t, []string{
				deepseek.DeepSeekChat,
				deepseek.DeepSeekCoder,
				deepseek.DeepSeekReasoner,
			}, model.ID)
		}
	}
}
