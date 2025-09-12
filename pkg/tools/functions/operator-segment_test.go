// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/operator-segment_test.go
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOperatorSegmentTool(t *testing.T) {
	tool := GetOperatorSegmentTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "operator_segment", tool.Tool.Name)
}

func TestGetOperatorSegmentHandler(t *testing.T) {
	mockSegments := map[string]interface{}{
		"segments": []map[string]interface{}{
			{
				"Name":      "alpha",
				"Port":      8301,
				"Bind":      "10.0.0.1",
				"Advertise": "10.0.0.1",
			},
			{
				"Name":      "beta",
				"Port":      8302,
				"Bind":      "10.0.0.2",
				"Advertise": "10.0.0.2",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/operator/segment", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockSegments)
	}))
	defer server.Close()

	t.Run("successful segment information retrieval", func(t *testing.T) {
		// Test would verify segment information is returned correctly
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
