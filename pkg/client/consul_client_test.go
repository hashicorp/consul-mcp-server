// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConsulClientFromContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		expectError bool
		errorMsg    string
	}{
		{
			name:        "missing session_id",
			ctx:         context.Background(),
			expectError: true,
			errorMsg:    "session_id not found in context",
		},
		{
			name:        "empty session_id",
			ctx:         context.WithValue(context.Background(), "session_id", ""),
			expectError: true,
			errorMsg:    "session_id not found in context",
		},
		{
			name:        "valid session_id",
			ctx:         context.WithValue(context.Background(), "session_id", "test-session"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewConsulClientFromContext(tt.ctx, logrus.New())

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestNewConsulClient(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("CONSUL_HTTP_ADDR", "http://test.consul:8500")
	os.Setenv("CONSUL_HTTP_TOKEN", "test-token")
	defer func() {
		os.Unsetenv("CONSUL_HTTP_ADDR")
		os.Unsetenv("CONSUL_HTTP_TOKEN")
	}()

	ctx := context.Background()
	sessionId := "test-session-123"

	client := NewConsulClient(ctx, sessionId, logrus.New())

	assert.NotNil(t, client)
	assert.Equal(t, sessionId, client.SessionID)
	assert.Equal(t, "http://test.consul:8500", client.Address)
	assert.Equal(t, "test-token", client.Token)
	assert.NotNil(t, client.client)

	// Test that client is stored in activeHttpClients
	storedClient := GetConsulHttpClient(sessionId)
	assert.Equal(t, client, storedClient)
}

func TestNewConsulClientWithContextOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("CONSUL_HTTP_ADDR", "http://env.consul:8500")
	defer os.Unsetenv("CONSUL_HTTP_ADDR")

	// Create context with overrides
	ctx := context.WithValue(context.Background(), "consul_address", "http://ctx.consul:8500")
	sessionId := "test-session-ctx"

	client := NewConsulClient(ctx, sessionId, logrus.New())

	assert.Equal(t, "http://ctx.consul:8500", client.Address)
}

func TestConsulClientHTTPMethods(t *testing.T) {
	// Create test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		assert.Contains(t, r.Header.Get("User-Agent"), "consul-mcp-server/")
		assert.Equal(t, "test-token", r.Header.Get("X-Consul-Token"))

		response := map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
		}

		if r.Method == "POST" || r.Method == "PUT" {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			// Read and echo back the request body
			body, _ := json.Marshal(response)
			w.Write(body)
		} else {
			body, _ := json.Marshal(response)
			w.Write(body)
		}
	}))
	defer testServer.Close()

	client := &ConsulHttpClient{
		SessionID: "test-session",
		Address:   testServer.URL,
		Token:     "test-token",
		client:    &http.Client{},
		Logger:    logrus.New(),
	}

	t.Run("GET request", func(t *testing.T) {
		resp, err := client.Get("test/endpoint")
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(resp, &result)
		require.NoError(t, err)
		assert.Equal(t, "GET", result["method"])
		assert.Equal(t, "/v1/test/endpoint", result["path"])
	})

	t.Run("POST request", func(t *testing.T) {
		data := map[string]string{"key": "value"}
		resp, err := client.Post("test/endpoint", data)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(resp, &result)
		require.NoError(t, err)
		assert.Equal(t, "POST", result["method"])
	})

	t.Run("PUT request", func(t *testing.T) {
		data := map[string]string{"key": "value"}
		resp, err := client.Put("test/endpoint", data)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(resp, &result)
		require.NoError(t, err)
		assert.Equal(t, "PUT", result["method"])
	})
}

func TestGetGetConsulHttpClientFromContext(t *testing.T) {
	t.Run("no session in context", func(t *testing.T) {
		ctx := context.Background()
		client, err := GetGetConsulHttpClientFromContext(ctx, logrus.New())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no active session")
		assert.Nil(t, client)
	})
}

type clientSessionKey struct{}

func TestDeleteConsulHttpClientForSession(t *testing.T) {
	sessionId := "delete-test-session"

	// Create a client
	ctx := context.WithValue(context.Background(), "session_id", sessionId)
	client := NewConsulClient(ctx, sessionId, logrus.New())

	// Verify it exists
	assert.Equal(t, client, GetConsulHttpClient(sessionId))

	// Delete it
	DeleteConsulHttpClientForSession(sessionId)

	// Verify it's gone
	assert.Nil(t, GetConsulHttpClient(sessionId))
}

func TestConsulClientErrorHandling(t *testing.T) {
	// Create test server that returns 404
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer testServer.Close()

	client := &ConsulHttpClient{
		SessionID: "test-session",
		Address:   testServer.URL,
		Token:     "test-token",
		client:    &http.Client{},
		Logger:    logrus.New(),
	}

	_, err := client.Get("test/endpoint")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}

func TestParseSkipTLSVerify(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected bool
	}{
		{
			name:     "no value in context",
			ctx:      context.Background(),
			expected: false,
		},
		{
			name:     "true value",
			ctx:      context.WithValue(context.Background(), contextKey(ConsulSkipTLSVerify), "true"),
			expected: true,
		},
		{
			name:     "false value",
			ctx:      context.WithValue(context.Background(), contextKey(ConsulSkipTLSVerify), "false"),
			expected: false,
		},
		{
			name:     "invalid value",
			ctx:      context.WithValue(context.Background(), contextKey(ConsulSkipTLSVerify), "invalid"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSkipTLSVerify(tt.ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}
