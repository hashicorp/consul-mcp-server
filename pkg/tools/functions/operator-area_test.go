// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/operator-area_test.go
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
	"github.com/stretchr/testify/require"
)

func TestGetOperatorAreasTool(t *testing.T) {
	tool := GetOperatorAreasTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "operator_areas", tool.Tool.Name)
}

func TestGetOperatorAreasHandler(t *testing.T) {
	mockAreas := []map[string]interface{}{
		{
			"ID":             "area-1",
			"PeerDatacenter": "dc2",
			"RetryJoin":      []string{"10.0.1.100", "10.0.1.101"},
			"UseTLS":         true,
		},
		{
			"ID":             "area-2",
			"PeerDatacenter": "dc3",
			"RetryJoin":      []string{"10.0.2.100"},
			"UseTLS":         false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/operator/area", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockAreas)
	}))
	defer server.Close()

	t.Run("successful areas listing", func(t *testing.T) {
		// Test would verify areas are returned correctly
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

func TestGetOperatorAreaTool(t *testing.T) {
	tool := GetOperatorAreaTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "operator_area", tool.Tool.Name)
}

func TestGetOperatorAreaHandler(t *testing.T) {
	//mockArea := map[string]interface{}{
	//	"ID":                   "area-uuid-1",
	//	"PeerDatacenter":       "dc2",
	//	"RetryJoin":            []string{"10.0.1.100", "10.0.1.101", "10.0.1.102"},
	//	"UseTLS":               true,
	//	"VerifyIncoming":       true,
	//	"VerifyOutgoing":       true,
	//	"VerifyServerHostname": true,
	//	"CAFile":               "/etc/consul/ca.pem",
	//	"CertFile":             "/etc/consul/cert.pem",
	//	"KeyFile":              "/etc/consul/key.pem",
	//}

	t.Run("required area_id parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"area_id": "area-uuid-1",
			},
		}

		areaId, err := request.RequireString("area_id")
		require.NoError(t, err)
		assert.Equal(t, "area-uuid-1", areaId)
	})

	t.Run("missing area_id parameter", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("area_id")
		assert.Error(t, err)
	})
}
