package base

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		envBase string
		url     string
		opts    []ClientOption
		wantErr bool
	}{
		{
			name:    "valid client with URL",
			envBase: "TEST",
			url:     "https://example.com",
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "valid client with options",
			envBase: "TEST",
			url:     "https://example.com",
			opts: []ClientOption{
				WithTimeout(60 * time.Second),
				WithRetries(3),
				WithVendor("TestVendor"),
				WithProduct("TestProduct"),
				WithBuild("1.0.0"),
			},
			wantErr: false,
		},
		{
			name:    "empty URL",
			envBase: "TEST",
			url:     "",
			opts:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := append([]ClientOption{WithURL(tt.url)}, tt.opts...)
			client, err := NewClient(tt.envBase, "", opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client, err := NewClient("TEST", server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	var result struct {
		Status string `json:"status"`
	}

	_, err = client.Get(ctx, "/test", &result)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result.Status)
	}
}

func TestClientPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"received": body["test"],
		})
	}))
	defer server.Close()

	client, err := NewClient("TEST", server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	body := map[string]string{"test": "value"}
	var result struct {
		Received string `json:"received"`
	}

	_, err = client.Post(ctx, "/test", body, &result)
	if err != nil {
		t.Fatalf("Post() error = %v", err)
	}

	if result.Received != "value" {
		t.Errorf("Expected received 'value', got '%s'", result.Received)
	}
}

func TestClientErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
	}))
	defer server.Close()

	client, err := NewClient("TEST", server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	var result interface{}

	_, err = client.Get(ctx, "/notfound", &result)
	if err == nil {
		t.Error("Expected error for 404 response")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", apiErr.StatusCode)
	}

	if !apiErr.IsNotFound() {
		t.Error("Expected IsNotFound() to return true")
	}
}

func TestAPIKeyAuth(t *testing.T) {
	var receivedHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("X-APIKeys")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client, err := NewClient("TEST", server.URL,
		WithAPIKeys("testaccess", "testsecret"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	var result interface{}
	_, _ = client.Get(ctx, "/test", &result)

	expected := "accessKey=testaccess; secretKey=testsecret"
	if receivedHeader != expected {
		t.Errorf("Expected header '%s', got '%s'", expected, receivedHeader)
	}
}

