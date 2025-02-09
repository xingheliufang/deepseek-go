package deepseek_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cohesion-org/deepseek-go"
)

func TestClientCreateChatCompletion(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected path /chat/completions, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer token" {
			t.Errorf("expected auth header Bearer token, got %s", authHeader)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"test-id"}`))
	}))
	defer testServer.Close()

	client, err := deepseek.NewClientWithOptions("token", deepseek.WithBaseURL(testServer.URL+"/"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.CreateChatCompletion(context.Background(), &deepseek.ChatCompletionRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "test-id" {
		t.Errorf("expected ID test-id, got %s", resp.ID)
	}
}

func TestCreateChatCompletion_NilRequest(t *testing.T) {
	client, err := deepseek.NewClientWithOptions("token")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CreateChatCompletion(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("expected error 'request cannot be nil', got %q", err.Error())
	}
}

func TestCreateChatCompletion_ErrorHandling(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"message": "invalid request"}}`))
	}))
	defer testServer.Close()

	client, err := deepseek.NewClientWithOptions("token", deepseek.WithBaseURL(testServer.URL+"/"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CreateChatCompletion(context.Background(), &deepseek.ChatCompletionRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*deepseek.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Message != "Bad request" {
		t.Errorf("expected error message 'Bad request', got %s", apiErr.Message)
	}
}

func TestCreateChatCompletionStream_ErrorHandling(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"message": "stream error"}}`))
	}))
	defer testServer.Close()

	client, err := deepseek.NewClientWithOptions("token", deepseek.WithBaseURL(testServer.URL+"/"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CreateChatCompletionStream(context.Background(), &deepseek.StreamChatCompletionRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*deepseek.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Message != "Bad request" {
		t.Errorf("expected error message 'Bad request', got %s", apiErr.Message)
	}
}
