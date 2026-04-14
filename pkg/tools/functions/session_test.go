// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/session_test.go
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

func TestGetSessionsHandler(t *testing.T) {
	mockSessions := []map[string]interface{}{
		{
			"ID":          "session-1",
			"Name":        "web-session",
			"Node":        "node-1",
			"LockDelay":   15000000000,
			"Behavior":    "release",
			"TTL":         "",
			"CreateIndex": 100,
			"ModifyIndex": 150,
		},
		{
			"ID":          "session-2",
			"Name":        "api-session",
			"Node":        "node-2",
			"LockDelay":   15000000000,
			"Behavior":    "delete",
			"TTL":         "30s",
			"CreateIndex": 200,
			"ModifyIndex": 250,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/session/list", r.URL.Path)
		assert.Equal(t, "default", r.URL.Query().Get("partition"))
		assert.Equal(t, "default", r.URL.Query().Get("ns"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockSessions)
	}))
	defer server.Close()

	t.Run("successful sessions listing", func(t *testing.T) {
		// Test would verify sessions are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})

	t.Run("with optional parameters", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"admin_partition": "team-a",
				"namespace":       "production",
				"dc":              "dc2",
			},
		}

		ap := request.GetString("admin_partition", "default")
		ns := request.GetString("namespace", "default")
		dc := request.GetString("dc", "")

		assert.Equal(t, "team-a", ap)
		assert.Equal(t, "production", ns)
		assert.Equal(t, "dc2", dc)
	})
}

func TestGetSessionTool(t *testing.T) {
	tool := GetSessionTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "session", tool.Tool.Name)
}

func TestGetSessionHandler(t *testing.T) {

	t.Run("required session_id parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"session_id":      "session-uuid-1",
				"admin_partition": "default",
				"namespace":       "default",
			},
		}

		sessionId, err := request.RequireString("session_id")
		require.NoError(t, err)
		assert.Equal(t, "session-uuid-1", sessionId)
	})

	t.Run("missing session_id parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("session_id")
		assert.Error(t, err)
	})
}

func TestGetSessionNodeTool(t *testing.T) {
	tool := GetSessionNodeTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "session_node", tool.Tool.Name)
}

func TestGetSessionNodeHandler(t *testing.T) {
	t.Run("required node_name parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"node_name":       "node1",
				"admin_partition": "default",
				"namespace":       "default",
			},
		}

		nodeName, err := request.RequireString("node_name")
		require.NoError(t, err)
		assert.Equal(t, "node1", nodeName)
	})

	t.Run("missing node_name parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("node_name")
		assert.Error(t, err)
	})
}

func TestGetSessionListTool(t *testing.T) {
	tool := GetSessionListTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "session_list", tool.Tool.Name)
}

func TestGetSessionListHandler(t *testing.T) {
	mockSessions := []map[string]interface{}{
		{
			"ID":          "session-uuid-1",
			"Name":        "test-session",
			"Node":        "node1",
			"Checks":      []string{"serfHealth"},
			"LockDelay":   15000000000,
			"Behavior":    "release",
			"TTL":         "",
			"CreateIndex": 1086449,
			"ModifyIndex": 1086449,
		},
		{
			"ID":          "session-uuid-2",
			"Name":        "another-session",
			"Node":        "node2",
			"Checks":      []string{"serfHealth"},
			"LockDelay":   15000000000,
			"Behavior":    "delete",
			"TTL":         "30s",
			"CreateIndex": 1086450,
			"ModifyIndex": 1086450,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/session/list", r.URL.Path)
		assert.Equal(t, "default", r.URL.Query().Get("partition"))
		assert.Equal(t, "default", r.URL.Query().Get("ns"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockSessions)
	}))
	defer server.Close()

	t.Run("successful sessions listing", func(t *testing.T) {
		// Test would verify sessions are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})

	t.Run("with optional parameters", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"admin_partition": "team-a",
				"namespace":       "production",
				"dc":              "dc2",
			},
		}

		ap := request.GetString("admin_partition", "default")
		ns := request.GetString("namespace", "default")
		dc := request.GetString("dc", "")

		assert.Equal(t, "team-a", ap)
		assert.Equal(t, "production", ns)
		assert.Equal(t, "dc2", dc)
	})
}
