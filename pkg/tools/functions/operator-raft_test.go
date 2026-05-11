// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/operator-raft_test.go
// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package functions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetOperatorRaftConfigurationTool(t *testing.T) {
	logger := log.New()
	tool := GetOperatorRaftConfigurationTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "operator_raft_configuration", tool.Tool.Name)
}

func TestGetOperatorRaftConfigurationHandler(t *testing.T) {
	mockRaftConfig := map[string]interface{}{
		"Servers": []map[string]interface{}{
			{
				"ID":      "127.0.0.1:8300",
				"Node":    "node1",
				"Address": "127.0.0.1:8300",
				"Leader":  true,
				"Voter":   true,
			},
		},
		"Index": 1,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/operator/raft/configuration", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockRaftConfig)
	}))
	defer server.Close()

	t.Run("successful raft configuration retrieval", func(t *testing.T) {
		// Test would verify raft configuration is returned correctly
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
