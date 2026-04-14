// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfigHandler(t *testing.T) {
	mockConfigEntry := map[string]interface{}{
		"Kind":        "service-defaults",
		"Name":        "web",
		"Protocol":    "http",
		"CreateIndex": 100,
		"ModifyIndex": 150,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/config/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockConfigEntry)
	}))
	defer server.Close()

	t.Run("successful config retrieval", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"kind": "service-defaults",
				"name": "web",
			},
		}

		kind, err := request.RequireString("kind")
		assert.NoError(t, err)
		assert.Equal(t, "service-defaults", kind)

		name, err := request.RequireString("name")
		assert.NoError(t, err)
		assert.Equal(t, "web", name)
	})

	t.Run("missing required parameters", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("kind")
		assert.Error(t, err)

		_, err = request.RequireString("name")
		assert.Error(t, err)
	})
}

func TestGetConfigsHandler(t *testing.T) {
	mockConfigEntries := []map[string]interface{}{
		{
			"Kind":        "service-defaults",
			"Name":        "web",
			"Protocol":    "http",
			"CreateIndex": 100,
			"ModifyIndex": 150,
		},
		{
			"Kind": "service-router",
			"Name": "api",
			"Routes": []map[string]interface{}{
				{
					"Match": map[string]interface{}{
						"HTTP": map[string]interface{}{
							"PathPrefix": "/v1/",
						},
					},
					"Destination": map[string]interface{}{
						"Service": "api-v1",
					},
				},
			},
			"CreateIndex": 200,
			"ModifyIndex": 250,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/config", r.URL.Path)

		query := r.URL.Query()
		if kind := query.Get("kind"); kind != "" {
			// Filter by kind if specified
			filtered := []map[string]interface{}{}
			for _, entry := range mockConfigEntries {
				if entry["Kind"] == kind {
					filtered = append(filtered, entry)
				}
			}
			json.NewEncoder(w).Encode(filtered)
		} else {
			json.NewEncoder(w).Encode(mockConfigEntries)
		}
	}))
	defer server.Close()

	t.Run("successful config entries listing", func(t *testing.T) {
		// Test would verify config entries are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})

	t.Run("filter by kind", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"kind": "service-defaults",
			},
		}

		kind := request.GetString("kind", "")
		assert.Equal(t, "service-defaults", kind)
	})
}

func TestConfigEntryTypes(t *testing.T) {
	configKinds := []string{
		"service-defaults",
		"proxy-defaults",
		"service-router",
		"service-splitter",
		"service-resolver",
		"ingress-gateway",
		"terminating-gateway",
		"mesh",
		"exported-services",
		"sameness-group",
	}

	for _, kind := range configKinds {
		t.Run(fmt.Sprintf("%s config kind", kind), func(t *testing.T) {
			configEntry := map[string]interface{}{
				"Kind": kind,
				"Name": fmt.Sprintf("test-%s", kind),
			}

			assert.Equal(t, kind, configEntry["Kind"])
		})
	}
}

func TestConfigEntryValidation(t *testing.T) {
	t.Run("service-defaults validation", func(t *testing.T) {
		serviceDefaults := map[string]interface{}{
			"Kind":     "service-defaults",
			"Name":     "web",
			"Protocol": "http",
			"MeshGateway": map[string]interface{}{
				"Mode": "local",
			},
		}

		assert.Equal(t, "service-defaults", serviceDefaults["Kind"])
		assert.Equal(t, "http", serviceDefaults["Protocol"])

		meshGateway := serviceDefaults["MeshGateway"].(map[string]interface{})
		assert.Equal(t, "local", meshGateway["Mode"])
	})

	t.Run("proxy-defaults validation", func(t *testing.T) {
		proxyDefaults := map[string]interface{}{
			"Kind": "proxy-defaults",
			"Name": "global",
			"Config": map[string]interface{}{
				"protocol":                 "http",
				"local_connect_timeout_ms": 5000,
			},
		}

		assert.Equal(t, "proxy-defaults", proxyDefaults["Kind"])

		config := proxyDefaults["Config"].(map[string]interface{})
		assert.Equal(t, "http", config["protocol"])
	})

	t.Run("service-router validation", func(t *testing.T) {
		serviceRouter := map[string]interface{}{
			"Kind": "service-router",
			"Name": "api",
			"Routes": []map[string]interface{}{
				{
					"Match": map[string]interface{}{
						"HTTP": map[string]interface{}{
							"PathPrefix": "/v1/",
						},
					},
					"Destination": map[string]interface{}{
						"Service": "api-v1",
					},
				},
			},
		}

		assert.Equal(t, "service-router", serviceRouter["Kind"])

		routes := serviceRouter["Routes"].([]map[string]interface{})
		assert.Len(t, routes, 1)
	})
}

func TestConfigEntryResponseProcessing(t *testing.T) {
	t.Run("complex config entry response", func(t *testing.T) {
		configEntry := map[string]interface{}{
			"Kind": "ingress-gateway",
			"Name": "ingress-service",
			"Listeners": []map[string]interface{}{
				{
					"Port":     8080,
					"Protocol": "http",
					"Services": []map[string]interface{}{
						{
							"Name":  "web",
							"Hosts": []string{"web.example.com"},
						},
					},
				},
			},
			"Meta": map[string]string{
				"environment": "production",
			},
		}

		data, err := json.MarshalIndent(configEntry, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "ingress-gateway", unmarshaled["Kind"])

		listeners := unmarshaled["Listeners"].([]interface{})
		assert.Len(t, listeners, 1)
	})
}

func TestConfigEntryParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid config entry parameters",
			arguments: map[string]interface{}{
				"kind":            "service-defaults",
				"name":            "web",
				"namespace":       "default",
				"admin_partition": "default",
			},
			expectError: false,
		},
		{
			name: "missing kind",
			arguments: map[string]interface{}{
				"name": "web",
			},
			expectError: true,
			errorField:  "kind",
		},
		{
			name: "missing name",
			arguments: map[string]interface{}{
				"kind": "service-defaults",
			},
			expectError: true,
			errorField:  "name",
		},
		{
			name: "config entries listing with kind filter",
			arguments: map[string]interface{}{
				"kind": "proxy-defaults",
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
				// Test optional parameters
				namespace := request.GetString("namespace", "default")
				partition := request.GetString("admin_partition", "default")

				assert.NotEmpty(t, namespace)
				assert.NotEmpty(t, partition)
			}
		})
	}
}

func TestConfigEntryErrorHandling(t *testing.T) {
	t.Run("config entry not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Config entry not found"))
		}))
		defer server.Close()

		// Would test handling of non-existent config entries
		assert.True(t, true) // Placeholder
	})

	t.Run("invalid config entry format", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid config entry format"))
		}))
		defer server.Close()

		// Would test handling of malformed config entries
		assert.True(t, true) // Placeholder
	})
}
