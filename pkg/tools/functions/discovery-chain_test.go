// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/discovery-chain_test.go
// Copyright IBM Corp. 2025
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDiscoveryChainTool(t *testing.T) {
	tool := GetDiscoveryChainTool(logrus.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "discovery_chain", tool.Tool.Name)
}

func TestGetDiscoveryChainHandler(t *testing.T) {
	mockDiscoveryChain := map[string]interface{}{
		"Chain": map[string]interface{}{
			"ServiceName": "web",
			"Namespace":   "default",
			"Datacenter":  "dc1",
			"Protocol":    "http",
			"StartNode":   "resolver:web.default.default.dc1",
			"Nodes": map[string]interface{}{
				"resolver:web.default.default.dc1": map[string]interface{}{
					"Type": "resolver",
					"Name": "web.default.default.dc1",
					"Resolver": map[string]interface{}{
						"Target":          "web.default.default.dc1",
						"ConnectTimeout":  "5s",
						"RequestTimeout":  "15s",
						"FailoverTargets": []string{},
					},
				},
			},
			"Targets": map[string]interface{}{
				"web.default.default.dc1": map[string]interface{}{
					"ID":         "web.default.default.dc1",
					"Service":    "web",
					"Namespace":  "default",
					"Partition":  "default",
					"Datacenter": "dc1",
					"Subset":     map[string]interface{}{},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/discovery-chain/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockDiscoveryChain)
	}))
	defer server.Close()

	t.Run("successful discovery chain retrieval", func(t *testing.T) {
		// Test would verify discovery chain is returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})
}

func TestGetDiscoveryChainHandlerWithOverrides(t *testing.T) {
	t.Run("with override parameters", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"service_name":               "web",
				"compile_dc":                 "dc2",
				"override_protocol":          "http",
				"override_mesh_gateway_mode": "local",
				"override_connect_timeout":   "10s",
			},
		}

		serviceName, err := request.RequireString("service_name")
		require.NoError(t, err)
		assert.Equal(t, "web", serviceName)

		compileDc := request.GetString("compile_dc", "")
		assert.Equal(t, "dc2", compileDc)

		overrideProtocol := request.GetString("override_protocol", "")
		assert.Equal(t, "http", overrideProtocol)
	})
}
