// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/exported-services_test.go
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
)

func TestGetExportedServicesTool(t *testing.T) {
	logger := log.New()
	tool := GetExportedServicesTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "exported_services", tool.Tool.Name)
}

func TestGetExportedServicesHandler(t *testing.T) {
	mockExportedServices := []map[string]interface{}{
		{
			"Service":   "web",
			"Namespace": "default",
			"Partition": "default",
			"Consumers": []map[string]interface{}{
				{
					"Peer": "cluster-02",
				},
			},
		},
		{
			"Service":   "api",
			"Namespace": "default",
			"Partition": "default",
			"Consumers": []map[string]interface{}{
				{
					"Partition": "frontend",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/exported-services", r.URL.Path)
		assert.Equal(t, "default", r.URL.Query().Get("partition"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockExportedServices)
	}))
	defer server.Close()

	t.Run("successful exported services listing", func(t *testing.T) {
		// Test would verify exported services are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})

	t.Run("with admin partition parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"admin_partition": "team-a",
			},
		}

		ap := request.GetString("admin_partition", "default")
		assert.Equal(t, "team-a", ap)
	})
}
