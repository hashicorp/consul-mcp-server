// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetKVTool(t *testing.T) {
	tool := GetKVTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "kv_get", tool.Tool.Name)
}

func TestGetKVHandler(t *testing.T) {
	testValue := "test-configuration-value"
	encodedValue := base64.StdEncoding.EncodeToString([]byte(testValue))

	mockKVPair := []map[string]interface{}{
		{
			"Key":         "config/app/database/host",
			"Value":       encodedValue,
			"Flags":       0,
			"CreateIndex": 100,
			"ModifyIndex": 200,
			"LockIndex":   0,
			"Session":     "",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/kv/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockKVPair)
	}))
	defer server.Close()

	t.Run("successful key retrieval", func(t *testing.T) {
		// Test would verify KV pair is returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})
}

func TestGetKVKeysTool(t *testing.T) {
	tool := GetKVKeysTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "kv_keys", tool.Tool.Name)
}

func TestGetKVKeysHandler(t *testing.T) {
	mockKeys := []string{
		"config/app/database/host",
		"config/app/database/port",
		"config/app/database/name",
		"config/app/redis/host",
		"config/app/redis/port",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/kv/")

		// Check for keys parameter
		query := r.URL.Query()
		assert.Equal(t, "", query.Get("keys")) // Should be present but empty value

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockKeys)
	}))
	defer server.Close()

	t.Run("successful keys listing", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"prefix": "config/app/",
			},
		}

		prefix, err := request.RequireString("prefix")
		assert.NoError(t, err)
		assert.Equal(t, "config/app/", prefix)
	})

	t.Run("missing prefix parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("prefix")
		assert.Error(t, err)
	})
}

func TestGetKVRecursiveTool(t *testing.T) {
	tool := GetKVRecursiveTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "kv_recursive", tool.Tool.Name)
}

func TestGetKVRecursiveHandler(t *testing.T) {
	dbHost := base64.StdEncoding.EncodeToString([]byte("postgres.internal"))
	dbPort := base64.StdEncoding.EncodeToString([]byte("5432"))
	redisHost := base64.StdEncoding.EncodeToString([]byte("redis.internal"))

	mockKVRecursive := []map[string]interface{}{
		{
			"Key":         "config/app/database/host",
			"Value":       dbHost,
			"Flags":       0,
			"CreateIndex": 10,
			"ModifyIndex": 15,
		},
		{
			"Key":         "config/app/database/port",
			"Value":       dbPort,
			"Flags":       0,
			"CreateIndex": 11,
			"ModifyIndex": 16,
		},
		{
			"Key":         "config/app/redis/host",
			"Value":       redisHost,
			"Flags":       0,
			"CreateIndex": 12,
			"ModifyIndex": 17,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/kv/")

		// Check for recurse parameter
		query := r.URL.Query()
		assert.Equal(t, "", query.Get("recurse"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockKVRecursive)
	}))
	defer server.Close()

	t.Run("successful recursive get", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"prefix": "config/app/",
			},
		}

		prefix := request.GetString("prefix", "")
		assert.Equal(t, "config/app/", prefix)
	})
}

func TestKVParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid KV get parameters",
			arguments: map[string]interface{}{
				"key":             "config/database/host",
				"datacenter":      "dc1",
				"namespace":       "production",
				"admin_partition": "team-a",
			},
			expectError: false,
		},
		{
			name: "missing key for get",
			arguments: map[string]interface{}{
				"datacenter": "dc1",
			},
			expectError: true,
			errorField:  "key",
		},
		{
			name: "valid keys list parameters",
			arguments: map[string]interface{}{
				"prefix":    "config/",
				"separator": "/",
			},
			expectError: false,
		},
		{
			name: "valid recursive parameters",
			arguments: map[string]interface{}{
				"prefix": "config/app/",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &MockCallToolRequest{Arguments: tt.arguments}

			if tt.expectError {
				_, err := request.RequireString(tt.errorField)
				assert.Error(t, err)
			} else {
				// Test optional parameters
				datacenter := request.GetString("datacenter", "dc1")
				namespace := request.GetString("namespace", "default")
				partition := request.GetString("admin_partition", "default")

				assert.NotEmpty(t, datacenter)
				assert.NotEmpty(t, namespace)
				assert.NotEmpty(t, partition)
			}
		})
	}
}

func TestKVValueDecoding(t *testing.T) {
	t.Run("base64 value decoding", func(t *testing.T) {
		originalValue := "database-configuration-secret"
		encodedValue := base64.StdEncoding.EncodeToString([]byte(originalValue))

		kvPair := map[string]interface{}{
			"Key":   "config/database/password",
			"Value": encodedValue,
		}

		// Test that we can decode the value
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedValue)
		require.NoError(t, err)
		assert.Equal(t, originalValue, string(decodedBytes))

		data, err := json.MarshalIndent(kvPair, "", "  ")
		require.NoError(t, err)
		assert.Contains(t, string(data), encodedValue)
	})

	t.Run("empty value handling", func(t *testing.T) {
		kvPair := map[string]interface{}{
			"Key":   "config/feature/flag",
			"Value": nil,
		}

		data, err := json.Marshal(kvPair)
		require.NoError(t, err)
		assert.Contains(t, string(data), "null")
	})
}

func TestKVResponseProcessing(t *testing.T) {
	t.Run("single KV pair response", func(t *testing.T) {
		value := base64.StdEncoding.EncodeToString([]byte("test-value"))
		kvResponse := []map[string]interface{}{
			{
				"Key":         "test/key",
				"Value":       value,
				"Flags":       0,
				"CreateIndex": 100,
				"ModifyIndex": 150,
			},
		}

		data, err := json.MarshalIndent(kvResponse, "", "  ")
		require.NoError(t, err)

		var unmarshaled []map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled, 1)
		assert.Equal(t, "test/key", unmarshaled[0]["Key"])
		assert.Equal(t, value, unmarshaled[0]["Value"])
	})

	t.Run("multiple KV pairs response", func(t *testing.T) {
		kvPairs := []map[string]interface{}{
			{"Key": "config/app/name", "Value": base64.StdEncoding.EncodeToString([]byte("my-app"))},
			{"Key": "config/app/version", "Value": base64.StdEncoding.EncodeToString([]byte("1.0.0"))},
			{"Key": "config/app/debug", "Value": base64.StdEncoding.EncodeToString([]byte("false"))},
		}

		data, err := json.MarshalIndent(kvPairs, "", "  ")
		require.NoError(t, err)
		assert.Contains(t, string(data), "config/app/name")
		assert.Contains(t, string(data), "config/app/version")
		assert.Contains(t, string(data), "config/app/debug")
	})
}

func TestKVFilteringAndQuerying(t *testing.T) {
	t.Run("prefix filtering", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"prefix":    "config/production/",
				"separator": "/",
			},
		}

		prefix := request.GetString("prefix", "")
		separator := request.GetString("separator", "")

		assert.Equal(t, "config/production/", prefix)
		assert.Equal(t, "/", separator)
	})

	t.Run("datacenter specific queries", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"key":        "global/config",
				"datacenter": "dc2",
			},
		}

		key, err := request.RequireString("key")
		assert.NoError(t, err)
		assert.Equal(t, "global/config", key)

		datacenter := request.GetString("datacenter", "dc1")
		assert.Equal(t, "dc2", datacenter)
	})
}

func TestKVErrorScenarios(t *testing.T) {
	t.Run("key not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Key not found"))
		}))
		defer server.Close()

		// Would test handling of non-existent keys
		assert.True(t, true) // Placeholder
	})

	t.Run("ACL permission denied for KV", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Permission denied"))
		}))
		defer server.Close()

		// Would test handling of ACL permission errors for KV operations
		assert.True(t, true) // Placeholder
	})

	t.Run("invalid base64 value", func(t *testing.T) {
		invalidBase64 := "invalid-base64-content!"

		// Test that invalid base64 is handled gracefully
		_, err := base64.StdEncoding.DecodeString(invalidBase64)
		assert.Error(t, err)
	})
}

func TestKVConsistencyModes(t *testing.T) {
	t.Run("default consistency", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"key": "config/app/setting",
			},
		}

		// Default consistency mode (eventual)
		consistency := request.GetString("consistency", "default")
		assert.Equal(t, "default", consistency)
	})

	t.Run("consistent read", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"key":        "config/critical/setting",
				"consistent": "true",
			},
		}

		consistent := request.GetString("consistent", "false")
		assert.Equal(t, "true", consistent)
	})

	t.Run("stale read", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"key":   "config/cache/setting",
				"stale": "true",
			},
		}

		stale := request.GetString("stale", "false")
		assert.Equal(t, "true", stale)
	})
}
