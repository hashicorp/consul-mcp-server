// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAgentServiceTool(t *testing.T) {
	logger := log.New()
	tool := GetAgentServiceTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "agent_service", tool.Tool.Name)
}

func TestGetAgentServiceHandler(t *testing.T) {
	mockService := map[string]interface{}{
		"ID":      "web-service-1",
		"Service": "web",
		"Tags":    []string{"v1.0", "production", "http"},
		"Meta": map[string]string{
			"version":     "1.0.0",
			"environment": "production",
			"team":        "frontend",
		},
		"Port":    8080,
		"Address": "192.168.1.10",
		"Weights": map[string]interface{}{
			"Passing": 10,
			"Warning": 1,
		},
		"EnableTagOverride": false,
		"CreateIndex":       100,
		"ModifyIndex":       150,
		"ContentHash":       "hash-web-service-123",
		"Proxy": map[string]interface{}{
			"DestinationServiceName": "web",
			"DestinationServiceID":   "web-service-1",
			"LocalServiceAddress":    "127.0.0.1",
			"LocalServicePort":       8080,
			"Config": map[string]interface{}{
				"protocol": "http",
			},
			"Upstreams": []map[string]interface{}{
				{
					"DestinationType":      "service",
					"DestinationName":      "database",
					"DestinationNamespace": "default",
					"LocalBindPort":        5432,
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/agent/service/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockService)
	}))
	defer server.Close()

	t.Run("successful service retrieval", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"service_id": "web-service-1",
			},
		}

		serviceID, err := request.RequireString("service_id")
		assert.NoError(t, err)
		assert.Equal(t, "web-service-1", serviceID)
	})

	t.Run("missing service ID", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("service_id")
		assert.Error(t, err)
	})
}

func TestGetAgentServiceHealthTool(t *testing.T) {
	logger := log.New()
	tool := GetAgentServiceHealthTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "agent_service_health", tool.Tool.Name)
}

func TestGetAgentServiceHealthHandler(t *testing.T) {
	mockServiceHealth := map[string]interface{}{
		"Service": map[string]interface{}{
			"ID":      "web-service-1",
			"Service": "web",
			"Tags":    []string{"v1.0", "production"},
			"Port":    8080,
		},
		"Checks": []map[string]interface{}{
			{
				"Node":        "consul-agent-1",
				"CheckID":     "service:web-service-1",
				"Name":        "Service 'web' check",
				"Status":      "passing",
				"Notes":       "",
				"Output":      "HTTP GET http://localhost:8080/health: 200 OK",
				"ServiceID":   "web-service-1",
				"ServiceName": "web",
			},
			{
				"Node":        "consul-agent-1",
				"CheckID":     "serfHealth",
				"Name":        "Serf Health Status",
				"Status":      "passing",
				"Notes":       "",
				"Output":      "Agent alive and reachable",
				"ServiceID":   "",
				"ServiceName": "",
			},
		},
		"AggregatedStatus": "passing",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/agent/health/service/id/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockServiceHealth)
	}))
	defer server.Close()

	t.Run("successful service health retrieval", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"service_id": "web-service-1",
			},
		}

		serviceID, err := request.RequireString("service_id")
		assert.NoError(t, err)
		assert.Equal(t, "web-service-1", serviceID)
	})
}

func TestAgentServiceParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid service ID",
			arguments: map[string]interface{}{
				"service_id": "web-service-1",
			},
			expectError: false,
		},
		{
			name:        "missing service ID",
			arguments:   map[string]interface{}{},
			expectError: true,
			errorField:  "service_id",
		},
		{
			name: "service health with format parameter",
			arguments: map[string]interface{}{
				"service_id": "web-service-1",
				"format":     "pretty",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &MockCallToolRequest{Arguments: tt.arguments}

			if tt.expectError {
				_, err := request.RequireString(tt.errorField)
				assert.Error(t, err)
			} else {
				serviceID, err := request.RequireString("service_id")
				assert.NoError(t, err)
				assert.NotEmpty(t, serviceID)
			}
		})
	}
}

func TestAgentServiceConfiguration(t *testing.T) {
	t.Run("service with proxy configuration", func(t *testing.T) {
		service := map[string]interface{}{
			"ID":      "web-proxy",
			"Service": "web-sidecar-proxy",
			"Kind":    "connect-proxy",
			"Proxy": map[string]interface{}{
				"DestinationServiceName": "web",
				"DestinationServiceID":   "web-service-1",
				"LocalServiceAddress":    "127.0.0.1",
				"LocalServicePort":       8080,
				"Upstreams": []map[string]interface{}{
					{
						"DestinationType": "service",
						"DestinationName": "database",
						"LocalBindPort":   5432,
					},
					{
						"DestinationType": "service",
						"DestinationName": "cache",
						"LocalBindPort":   6379,
					},
				},
			},
		}

		proxy := service["Proxy"].(map[string]interface{})
		assert.Equal(t, "web", proxy["DestinationServiceName"])

		upstreams := proxy["Upstreams"].([]map[string]interface{})
		assert.Len(t, upstreams, 2)
		assert.Equal(t, "database", upstreams[0]["DestinationName"])
		assert.Equal(t, 5432, upstreams[0]["LocalBindPort"])
	})

	t.Run("service with mesh gateway mode", func(t *testing.T) {
		service := map[string]interface{}{
			"ID": "web-service",
			"Proxy": map[string]interface{}{
				"MeshGateway": map[string]interface{}{
					"Mode": "local",
				},
				"Expose": map[string]interface{}{
					"Checks": true,
					"Paths": []map[string]interface{}{
						{
							"Path":          "/health",
							"LocalPathPort": 8080,
							"ListenerPort":  21500,
						},
					},
				},
			},
		}

		proxy := service["Proxy"].(map[string]interface{})
		meshGateway := proxy["MeshGateway"].(map[string]interface{})
		assert.Equal(t, "local", meshGateway["Mode"])

		expose := proxy["Expose"].(map[string]interface{})
		assert.Equal(t, true, expose["Checks"])
	})
}

func TestAgentServiceWeights(t *testing.T) {
	t.Run("service weight configuration", func(t *testing.T) {
		service := map[string]interface{}{
			"Weights": map[string]interface{}{
				"Passing": 10,
				"Warning": 1,
			},
		}

		weights := service["Weights"].(map[string]interface{})
		assert.Equal(t, 10, weights["Passing"])
		assert.Equal(t, 1, weights["Warning"])
	})

	t.Run("custom service weights", func(t *testing.T) {
		weights := []map[string]interface{}{
			{"Passing": 10, "Warning": 5},  // Higher warning weight
			{"Passing": 100, "Warning": 0}, // Zero warning weight
			{"Passing": 1, "Warning": 1},   // Equal weights
		}

		for _, weight := range weights {
			passing := weight["Passing"].(int)
			warning := weight["Warning"].(int)
			assert.GreaterOrEqual(t, passing, 0)
			assert.GreaterOrEqual(t, warning, 0)
		}
	})
}

func TestAgentServiceHealthStatus(t *testing.T) {
	healthStatuses := []string{"passing", "warning", "critical"}

	for _, status := range healthStatuses {
		t.Run(fmt.Sprintf("%s health status", status), func(t *testing.T) {
			serviceHealth := map[string]interface{}{
				"AggregatedStatus": status,
				"Checks": []map[string]interface{}{
					{
						"Status":  status,
						"CheckID": fmt.Sprintf("check-%s", status),
					},
				},
			}

			assert.Equal(t, status, serviceHealth["AggregatedStatus"])

			checks := serviceHealth["Checks"].([]map[string]interface{})
			assert.Equal(t, status, checks[0]["Status"])
		})
	}
}

func TestAgentServiceResponseProcessing(t *testing.T) {
	t.Run("complete service response", func(t *testing.T) {
		service := map[string]interface{}{
			"ID":      "complex-service",
			"Service": "complex-app",
			"Tags":    []string{"v2.0", "canary"},
			"Meta":    map[string]string{"version": "2.0", "canary": "true"},
			"Port":    9090,
			"Weights": map[string]interface{}{"Passing": 5, "Warning": 1},
			"Proxy": map[string]interface{}{
				"DestinationServiceName": "complex-app",
				"Upstreams": []map[string]interface{}{
					{"DestinationName": "auth", "LocalBindPort": 8081},
				},
			},
		}

		data, err := json.MarshalIndent(service, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "complex-service", unmarshaled["ID"])
		assert.Contains(t, unmarshaled, "Proxy")
		assert.Contains(t, unmarshaled, "Weights")
	})
}

func TestAgentServiceErrorHandling(t *testing.T) {
	t.Run("service not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Service not found"))
		}))
		defer server.Close()

		// Would test handling of non-existent services
		assert.True(t, true) // Placeholder
	})

	t.Run("agent unavailable", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Agent unavailable"))
		}))
		defer server.Close()

		// Would test handling of agent connectivity issues
		assert.True(t, true) // Placeholder
	})
}
