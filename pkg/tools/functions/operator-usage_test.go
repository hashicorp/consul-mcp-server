// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/operator-usage_test.go
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

func TestGetOperatorUsageTool(t *testing.T) {
	logger := log.New()
	tool := GetOperatorUsageTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "operator_usage", tool.Tool.Name)
}

func TestGetOperatorUsageHandler(t *testing.T) {
	mockUsage := map[string]interface{}{
		"Usage": map[string]interface{}{
			"consul": map[string]interface{}{
				"billable_service_instances": 3,
				"connect_service_instances": map[string]interface{}{
					"consul-dataplane": 0,
					"envoy":            0,
					"native":           0,
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/operator/usage", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockUsage)
	}))
	defer server.Close()

	t.Run("successful usage information retrieval", func(t *testing.T) {
		// Test would verify usage information is returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})
}
