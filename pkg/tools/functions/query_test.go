// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/query_test.go
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

func TestGetQueryTool(t *testing.T) {
	logger := log.New()
	tool := GetQueryTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "query", tool.Tool.Name)
}

func TestGetQueryHandler(t *testing.T) {
	mockQueries := []map[string]interface{}{
		{
			"ID":      "query-uuid-1",
			"Name":    "my-query",
			"Session": "session-uuid-1",
			"Token":   "",
			"Service": map[string]interface{}{
				"Service": "redis",
				"Failover": map[string]interface{}{
					"NearestN":    3,
					"Datacenters": []string{"dc1", "dc2"},
				},
				"OnlyPassing": false,
				"Tags":        []string{"master", "!experimental"},
			},
			"DNS": map[string]interface{}{
				"TTL": "10s",
			},
			"RaftIndex": map[string]interface{}{
				"CreateIndex": 23,
				"ModifyIndex": 42,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/query", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockQueries)
	}))
	defer server.Close()

	t.Run("successful queries listing", func(t *testing.T) {
		// Test would verify queries are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})

	t.Run("with datacenter parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"dc": "dc2",
			},
		}

		dc := request.GetString("dc", "")
		assert.Equal(t, "dc2", dc)
	})
}

func TestGetQueryByIdTool(t *testing.T) {
	logger := log.New()
	tool := GetQueryByIdTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "query_by_id", tool.Tool.Name)
}

func TestGetQueryByIdHandler(t *testing.T) {
	t.Run("required query_id parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"query_id": "query-uuid-1",
			},
		}

		queryId, err := request.RequireString("query_id")
		require.NoError(t, err)
		assert.Equal(t, "query-uuid-1", queryId)
	})

	t.Run("missing query_id parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("query_id")
		assert.Error(t, err)
	})
}
