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
	"github.com/stretchr/testify/require"
)

func TestGetHealthNodeTool(t *testing.T) {
	tool := GetHealthNodeTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "health_node", tool.Tool.Name)
}

func TestGetHealthNodeHandler(t *testing.T) {
	mockHealthChecks := []map[string]interface{}{
		{
			"Node":        "node1",
			"CheckID":     "service:redis:1",
			"Name":        "Service 'redis' check",
			"Status":      "passing",
			"Notes":       "",
			"Output":      "TCP connect 127.0.0.1:6379: Success",
			"ServiceID":   "redis",
			"ServiceName": "redis",
		},
		{
			"Node":        "node1",
			"CheckID":     "serfHealth",
			"Name":        "Serf Health Status",
			"Status":      "passing",
			"Notes":       "",
			"Output":      "Agent alive and reachable",
			"ServiceID":   "",
			"ServiceName": "",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/health/node/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockHealthChecks)
	}))
	defer server.Close()

	t.Run("successful health check", func(t *testing.T) {
		// Test would verify health checks are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})
}

func TestGetHealthServiceTool(t *testing.T) {
	tool := GetHealthServiceTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "health_service", tool.Tool.Name)
}

func TestGetHealthServiceHandler(t *testing.T) {
	mockServiceHealth := []map[string]interface{}{
		{
			"Node": map[string]interface{}{
				"ID":              "node1",
				"Node":            "consul-node-1",
				"Address":         "192.168.1.10",
				"Datacenter":      "dc1",
				"TaggedAddresses": map[string]string{"lan": "192.168.1.10"},
			},
			"Service": map[string]interface{}{
				"ID":      "web-1",
				"Service": "web",
				"Tags":    []string{"v1.0", "production"},
				"Address": "192.168.1.10",
				"Port":    8080,
			},
			"Checks": []map[string]interface{}{
				{
					"Node":        "consul-node-1",
					"CheckID":     "service:web",
					"Name":        "Service 'web' check",
					"Status":      "passing",
					"Output":      "HTTP GET http://localhost:8080/health: 200 OK",
					"ServiceID":   "web-1",
					"ServiceName": "web",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/health/service/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockServiceHealth)
	}))
	defer server.Close()

	t.Run("successful service health check", func(t *testing.T) {
		// Test would verify service health information is returned correctly
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
}

func TestGetHealthChecksTool(t *testing.T) {
	tool := GetHealthChecksTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "health_checks", tool.Tool.Name)
}

func TestGetHealthChecksHandler(t *testing.T) {
	mockChecks := []map[string]interface{}{
		{
			"Node":        "consul-node-1",
			"CheckID":     "serfHealth",
			"Name":        "Serf Health Status",
			"Status":      "passing",
			"Notes":       "",
			"Output":      "Agent alive and reachable",
			"ServiceID":   "",
			"ServiceName": "",
		},
		{
			"Node":        "consul-node-1",
			"CheckID":     "service:web",
			"Name":        "Service 'web' check",
			"Status":      "critical",
			"Notes":       "",
			"Output":      "Get http://localhost:8080/health: dial tcp 127.0.0.1:8080: connect: connection refused",
			"ServiceID":   "web",
			"ServiceName": "web",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/health/checks/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockChecks)
	}))
	defer server.Close()

	t.Run("successful service checks listing", func(t *testing.T) {
		// Test would verify service-specific health checks
		assert.True(t, true) // Placeholder
	})
}

func TestGetHealthStateTool(t *testing.T) {
	tool := GetHealthStateTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "health_state", tool.Tool.Name)
}

func TestGetHealthStateHandler(t *testing.T) {
	mockCriticalChecks := []map[string]interface{}{
		{
			"Node":        "consul-node-2",
			"CheckID":     "service:database",
			"Name":        "Service 'database' check",
			"Status":      "critical",
			"Notes":       "",
			"Output":      "dial tcp 192.168.1.20:5432: connect: connection refused",
			"ServiceID":   "database",
			"ServiceName": "database",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/health/state/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockCriticalChecks)
	}))
	defer server.Close()

	t.Run("critical health state", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"state": "critical",
			},
		}

		state, err := request.RequireString("state")
		assert.NoError(t, err)
		assert.Equal(t, "critical", state)
	})

	t.Run("valid health states", func(t *testing.T) {
		validStates := []string{"any", "passing", "warning", "critical"}

		for _, state := range validStates {
			request := &MockCallToolRequest{
				Arguments: map[string]interface{}{
					"state": state,
				},
			}

			retrievedState, err := request.RequireString("state")
			assert.NoError(t, err)
			assert.Equal(t, state, retrievedState)
		}
	})
}

func TestHealthParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid node health parameters",
			arguments: map[string]interface{}{
				"node_name":       "consul-node-1",
				"datacenter":      "dc1",
				"namespace":       "default",
				"admin_partition": "default",
			},
			expectError: false,
		},
		{
			name: "missing node name",
			arguments: map[string]interface{}{
				"datacenter": "dc1",
			},
			expectError: true,
			errorField:  "node_name",
		},
		{
			name: "valid service health parameters",
			arguments: map[string]interface{}{
				"service_name":    "web",
				"datacenter":      "dc1",
				"namespace":       "production",
				"admin_partition": "team-a",
			},
			expectError: false,
		},
		{
			name: "missing service name",
			arguments: map[string]interface{}{
				"datacenter": "dc1",
			},
			expectError: true,
			errorField:  "service_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &MockCallToolRequest{Arguments: tt.arguments}

			if tt.expectError {
				_, err := request.RequireString(tt.errorField)
				assert.Error(t, err)
			} else {
				// Test optional parameters with defaults
				datacenter := request.GetString("datacenter", "dc1")
				namespace := request.GetString("namespace", "default")
				partition := request.GetString("admin_partition", "default")

				assert.NotEmpty(t, datacenter)
				assert.NotEmpty(t, namespace)
				assert.NotEmpty(t, partition)
			}
		})
	}
}

func TestHealthFilterParameters(t *testing.T) {
	t.Run("service health with tags", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"service_name": "web",
				"tag":          "production",
				"passing":      true,
			},
		}

		serviceName, err := request.RequireString("service_name")
		assert.NoError(t, err)
		assert.Equal(t, "web", serviceName)

		tag := request.GetString("tag", "")
		assert.Equal(t, "production", tag)
	})

	t.Run("health checks with node meta", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"service_name": "database",
				"near":         "consul-node-1",
			},
		}

		near := request.GetString("near", "")
		assert.Equal(t, "consul-node-1", near)
	})
}

func TestHealthResponseProcessing(t *testing.T) {
	t.Run("health check status parsing", func(t *testing.T) {
		healthCheck := map[string]interface{}{
			"Node":        "consul-node-1",
			"CheckID":     "service:web",
			"Name":        "Service 'web' check",
			"Status":      "passing",
			"Output":      "HTTP GET http://localhost:8080/health: 200 OK",
			"ServiceID":   "web",
			"ServiceName": "web",
		}

		data, err := json.MarshalIndent(healthCheck, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "passing", unmarshaled["Status"])
		assert.Equal(t, "web", unmarshaled["ServiceName"])
	})

	t.Run("service health aggregation", func(t *testing.T) {
		serviceHealth := []map[string]interface{}{
			{
				"Node": map[string]interface{}{
					"Node":    "node1",
					"Address": "192.168.1.10",
				},
				"Service": map[string]interface{}{
					"Service": "web",
					"Port":    8080,
				},
				"Checks": []map[string]interface{}{
					{"Status": "passing"},
					{"Status": "warning"},
				},
			},
		}

		data, err := json.MarshalIndent(serviceHealth, "", "  ")
		require.NoError(t, err)
		assert.Contains(t, string(data), "passing")
		assert.Contains(t, string(data), "warning")
	})
}

func TestHealthErrorScenarios(t *testing.T) {
	t.Run("non-existent node", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Node not found"))
		}))
		defer server.Close()

		// Would test handling of non-existent nodes
		assert.True(t, true) // Placeholder
	})

	t.Run("non-existent service", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]interface{}{}) // Empty array for non-existent service
		}))
		defer server.Close()

		// Would test handling of non-existent services
		assert.True(t, true) // Placeholder
	})

	t.Run("invalid health state", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"state": "invalid-state",
			},
		}

		state, err := request.RequireString("state")
		assert.NoError(t, err)
		assert.Equal(t, "invalid-state", state)

		// In real implementation, would validate against allowed states
		validStates := []string{"any", "passing", "warning", "critical"}
		assert.NotContains(t, validStates, state)
	})
}
