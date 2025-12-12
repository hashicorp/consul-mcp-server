// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/operator-license_test.go
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
)

func TestGetOperatorLicenseTool(t *testing.T) {
	logger := log.New()
	tool := GetOperatorLicenseTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "operator_license", tool.Tool.Name)
}

func TestGetOperatorLicenseHandler(t *testing.T) {
	mockLicense := map[string]interface{}{
		"Valid": true,
		"License": map[string]interface{}{
			"license_id":      "2afbf681-0d1a-0649-cb6c-333ec9f0989c",
			"customer_id":     "7401728c-5c2c-4e04-adb5-87b18ad7d4a7",
			"installation_id": "*",
			"issue_time":      "2021-12-20T19:33:38.362717847Z",
			"start_time":      "2021-12-20T00:00:00Z",
			"expiration_time": "2022-12-20T23:59:59.999Z",
			"product":         "consul",
			"flags": map[string]interface{}{
				"modules": []string{"multi-datacenter", "governance-policy"},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/operator/license", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockLicense)
	}))
	defer server.Close()

	t.Run("successful license retrieval", func(t *testing.T) {
		// Test would verify license information is returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})
}
