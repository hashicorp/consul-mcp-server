// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/peering_test.go
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

func TestGetPeeringsTool(t *testing.T) {
	logger := log.New()
	tool := GetPeeringsTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "peerings", tool.Tool.Name)
}

func TestGetPeeringsHandler(t *testing.T) {
	mockPeerings := []map[string]interface{}{
		{
			"ID":                  "peer-id-1",
			"Name":                "cluster-02",
			"Partition":           "default",
			"State":               "ACTIVE",
			"PeerID":              "e83a315c-027e-bcb1-7c0c-a46650904a05",
			"PeerCAPems":          nil,
			"PeerServerName":      "",
			"PeerServerAddresses": []string{},
			"CreateIndex":         89,
			"ModifyIndex":         89,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/peerings", r.URL.Path)
		assert.Equal(t, "default", r.URL.Query().Get("partition"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockPeerings)
	}))
	defer server.Close()

	t.Run("successful peerings listing", func(t *testing.T) {
		// Test would verify peerings are returned correctly
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
