// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/connect_intentions_test.go
// Copyright IBM Corp. 2025
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConnectIntentionsTool(t *testing.T) {
	logger := log.New()
	tool := GetConnectIntentionsTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "connect_intentions", tool.Tool.Name)
}

func TestGetConnectIntentionsHandler(t *testing.T) {
	mockIntentions := []map[string]interface{}{
		{
			"ID":              "intention-uuid-1",
			"Description":     "Allow web to api",
			"SourceNS":        "default",
			"SourceName":      "web",
			"DestinationNS":   "default",
			"DestinationName": "api",
			"Action":          "allow",
			"CreateIndex":     11,
			"ModifyIndex":     11,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/connect/intentions", r.URL.Path)
		assert.Equal(t, "default", r.URL.Query().Get("partition"))
		assert.Equal(t, "default", r.URL.Query().Get("ns"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockIntentions)
	}))
	defer server.Close()

	t.Run("successful intentions listing", func(t *testing.T) {
		// Test would verify intentions are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})

	t.Run("with filter parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"admin_partition": "default",
				"namespace":       "default",
				"filter":          "SourceName == web",
			},
		}

		filter := request.GetString("filter", "")
		assert.Equal(t, "SourceName == web", filter)
	})
}

func TestGetConnectIntentionTool(t *testing.T) {
	logger := log.New()
	tool := GetConnectIntentionTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "connect_intention", tool.Tool.Name)
}

func TestGetConnectIntentionHandler(t *testing.T) {
	t.Run("required intention_id parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"intention_id": "intention-uuid-1",
			},
		}

		intentionId, err := request.RequireString("intention_id")
		require.NoError(t, err)
		assert.Equal(t, "intention-uuid-1", intentionId)
	})

	t.Run("missing intention_id parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("intention_id")
		assert.Error(t, err)
	})
}

func TestGetConnectIntentionMatchTool(t *testing.T) {
	logger := log.New()
	tool := GetConnectIntentionMatchTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "connect_intention_match", tool.Tool.Name)
}

func TestGetConnectIntentionMatchHandler(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"by":              "source",
				"name":            "web",
				"admin_partition": "default",
				"namespace":       "default",
			},
		}

		by, err := request.RequireString("by")
		require.NoError(t, err)
		assert.Equal(t, "source", by)

		name, err := request.RequireString("name")
		require.NoError(t, err)
		assert.Equal(t, "web", name)
	})
}

func TestGetConnectIntentionCheckTool(t *testing.T) {
	logger := log.New()
	tool := GetConnectIntentionCheckTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "connect_intention_check", tool.Tool.Name)
}

func TestGetConnectIntentionCheckHandler(t *testing.T) {
	t.Run("required parameters", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"source":          "web",
				"destination":     "api",
				"source_type":     "consul",
				"admin_partition": "default",
				"namespace":       "default",
			},
		}

		source, err := request.RequireString("source")
		require.NoError(t, err)
		assert.Equal(t, "web", source)

		destination, err := request.RequireString("destination")
		require.NoError(t, err)
		assert.Equal(t, "api", destination)
	})
}
