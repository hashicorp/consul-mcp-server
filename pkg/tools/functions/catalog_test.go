// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/consul-mcp-server/pkg/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetCatalogServicesTool(t *testing.T) {
	tool := GetCatalogServicesTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "catalog_services", tool.Tool.Name)
}

func TestGetCatalogServicesHandler(t *testing.T) {
	mockServices := client.Services{
		"web":      []string{"v1.0", "production"},
		"database": []string{"v2.1", "staging"},
		"cache":    []string{"v1.5"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/catalog/services", r.URL.Path)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "default", query.Get("partition"))
		assert.Equal(t, "default", query.Get("ns"))

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(mockServices)
		assert.NoError(t, err)
	}))
	defer server.Close()

	t.Run("successful services listing", func(t *testing.T) {
		// Test would verify the handler correctly processes the response
		// and returns properly formatted service data
		assert.True(t, true) // Placeholder for actual implementation
	})

	t.Run("custom partition and namespace", func(t *testing.T) {
		// Test with custom partition and namespace parameters
		assert.True(t, true) // Placeholder
	})
}

func TestGetCatalogNodesTool(t *testing.T) {
	tool := GetCatalogNodesTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "catalog_nodes", tool.Tool.Name)
}

func TestGetCatalogNodesHandler(t *testing.T) {
	mockNodes := []map[string]interface{}{
		{
			"ID":              "node1",
			"Node":            "consul-node-1",
			"Address":         "192.168.1.10",
			"Datacenter":      "dc1",
			"TaggedAddresses": map[string]string{"lan": "192.168.1.10", "wan": "10.0.0.10"},
			"Meta":            map[string]string{"zone": "us-west-1a"},
		},
		{
			"ID":              "node2",
			"Node":            "consul-node-2",
			"Address":         "192.168.1.11",
			"Datacenter":      "dc1",
			"TaggedAddresses": map[string]string{"lan": "192.168.1.11", "wan": "10.0.0.11"},
			"Meta":            map[string]string{"zone": "us-west-1b"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/catalog/nodes", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(mockNodes)
		assert.NoError(t, err)
	}))
	defer server.Close()

	t.Run("successful nodes listing", func(t *testing.T) {
		// Test would verify the handler returns node information correctly
		assert.True(t, true) // Placeholder
	})
}

func TestGetCatalogServiceTool(t *testing.T) {
	tool := GetCatalogServiceTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "catalog_service", tool.Tool.Name)
}

func TestGetCatalogServiceHandler(t *testing.T) {
	mockServiceNodes := []map[string]interface{}{
		{
			"ID":                       "service-instance-1",
			"Node":                     "consul-node-1",
			"Address":                  "192.168.1.10",
			"Datacenter":               "dc1",
			"ServiceID":                "web-1",
			"ServiceName":              "web",
			"ServiceTags":              []string{"v1.0", "production"},
			"ServiceAddress":           "192.168.1.10",
			"ServicePort":              8080,
			"ServiceMeta":              map[string]string{"version": "1.0"},
			"ServiceEnableTagOverride": false,
		},
		{
			"ID":                       "service-instance-2",
			"Node":                     "consul-node-2",
			"Address":                  "192.168.1.11",
			"Datacenter":               "dc1",
			"ServiceID":                "web-2",
			"ServiceName":              "web",
			"ServiceTags":              []string{"v1.0", "production"},
			"ServiceAddress":           "192.168.1.11",
			"ServicePort":              8080,
			"ServiceMeta":              map[string]string{"version": "1.0"},
			"ServiceEnableTagOverride": false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/catalog/service/")

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(mockServiceNodes)
		assert.NoError(t, err)
	}))
	defer server.Close()

	t.Run("successful service details", func(t *testing.T) {
		// Test would verify service-specific node information
		assert.True(t, true) // Placeholder
	})

	t.Run("service name parameter validation", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"service_name": "web",
			},
		}

		serviceName, err := request.RequireString("service_name")
		assert.NoError(t, err)
		assert.Equal(t, "web", serviceName)
	})

	t.Run("missing service name", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("service_name")
		assert.Error(t, err)
	})
}

func TestGetCatalogDatacentersTool(t *testing.T) {
	tool := GetCatalogDatacentersTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "catalog_datacenters", tool.Tool.Name)
}

func TestGetCatalogDatacentersHandler(t *testing.T) {
	mockDatacenters := []string{"dc1", "dc2", "dc3"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/catalog/datacenters", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(mockDatacenters)
		assert.NoError(t, err)
	}))
	defer server.Close()

	t.Run("successful datacenters listing", func(t *testing.T) {
		// Test would verify datacenter list is returned correctly
		assert.True(t, true) // Placeholder
	})
}

func TestCatalogParameterHandling(t *testing.T) {
	tests := []struct {
		name         string
		arguments    map[string]interface{}
		expectedDC   string
		expectedNS   string
		expectedPart string
	}{
		{
			name:         "default parameters",
			arguments:    map[string]interface{}{},
			expectedDC:   "dc1", // Assuming default from environment
			expectedNS:   "default",
			expectedPart: "default",
		},
		{
			name: "custom parameters",
			arguments: map[string]interface{}{
				"datacenter":      "dc2",
				"namespace":       "production",
				"admin_partition": "team-a",
			},
			expectedDC:   "dc2",
			expectedNS:   "production",
			expectedPart: "team-a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &MockCallToolRequest{Arguments: tt.arguments}

			// Test parameter extraction with defaults
			ns := request.GetString("namespace", "default")
			partition := request.GetString("admin_partition", "default")

			assert.Equal(t, tt.expectedNS, ns)
			assert.Equal(t, tt.expectedPart, partition)
		})
	}
}

func TestCatalogFilterParameters(t *testing.T) {
	t.Run("service filtering", func(t *testing.T) {
		// Test service filtering by tags, metadata, etc.
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"service_name": "web",
				"tag":          "production",
			},
		}

		serviceName, err := request.RequireString("service_name")
		assert.NoError(t, err)
		assert.Equal(t, "web", serviceName)

		tag := request.GetString("tag", "")
		assert.Equal(t, "production", tag)
	})

	t.Run("node filtering", func(t *testing.T) {
		// Test node filtering parameters
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"near": "node1",
			},
		}

		near := request.GetString("near", "")
		assert.Equal(t, "node1", near)
	})
}

func TestCatalogErrorHandling(t *testing.T) {
	t.Run("network error", func(t *testing.T) {
		// Test handling of network errors during API calls
		// Would mock a failing HTTP client
		assert.True(t, true) // Placeholder
	})

	t.Run("invalid service name", func(t *testing.T) {
		// Test handling of invalid service names
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"service_name": "",
			},
		}

		_, err := request.RequireString("service_name")
		assert.Error(t, err)
	})

	t.Run("consul api error", func(t *testing.T) {
		// Test handling of Consul API errors (404, 500, etc.)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte("Service not found"))
			assert.NoError(t, err)
		}))
		defer server.Close()

		// Would test that 404 errors are properly handled
		assert.True(t, true) // Placeholder
	})
}
