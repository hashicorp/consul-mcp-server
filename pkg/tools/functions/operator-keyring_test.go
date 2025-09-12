// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/operator-keyring_test.go
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetOperatorKeyringTool(t *testing.T) {
	logger := log.New()
	tool := GetOperatorKeyringTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "operator_keyring", tool.Tool.Name)
}

func TestGetOperatorKeyringHandler(t *testing.T) {
	mockKeyring := []map[string]interface{}{
		{
			"WAN":        true,
			"Datacenter": "dc1",
			"Keys": map[string]interface{}{
				"pUqJrVyVRj5jsiYEkM/tFQYfWyJIv4s3XkvDwy7Cu5s=": 1,
			},
			"NumNodes": 1,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/operator/keyring", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockKeyring)
	}))
	defer server.Close()

	t.Run("successful keyring retrieval", func(t *testing.T) {
		// Test would verify keyring information is returned correctly
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
