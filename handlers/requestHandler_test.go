package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cohesion-org/deepseek-go/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleSendChatCompletionRequest(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		resp, err := handlers.HandleSendChatCompletionRequest(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("request error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://invalid-url", nil)
		require.NoError(t, err)

		resp, err := handlers.HandleSendChatCompletionRequest(req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "error sending request:")
	})
}

func TestHandleNormalRequest(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		resp, err := handlers.HandelNormalRequest(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("request error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://invalid-url", nil)
		require.NoError(t, err)

		resp, err := handlers.HandelNormalRequest(req)
		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "error sending request:")
	})
}

func TestTimeoutConfiguration(t *testing.T) {
	t.Run("chat completion timeout", func(t *testing.T) {
		start := time.Now()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(250 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		resp, err := handlers.HandleSendChatCompletionRequest(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify timeout configuration isn't too short
		assert.WithinDuration(t, start, time.Now(), 300*time.Millisecond)
	})

	t.Run("normal request timeout", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		start := time.Now()
		resp, err := handlers.HandelNormalRequest(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify timeout allows successful completion
		assert.WithinDuration(t, start, time.Now(), 100*time.Millisecond)
	})
}

func TestErrorTimeout(t *testing.T) {
	t.Run("client timeout preservation", func(t *testing.T) {
		// Create test server that responds slower than client timeout
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(250 * time.Millisecond) // Longer than client timeout
		}))
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		// Create client with short timeout
		shortClient := &http.Client{
			Timeout: 5 * time.Millisecond,
		}

		_, err = shortClient.Do(req)
		require.Error(t, err)

		assert.True(t, errors.Is(err, context.DeadlineExceeded) ||
			(os.IsTimeout(err)), "should be timeout error")
	})
}
