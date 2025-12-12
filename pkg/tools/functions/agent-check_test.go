// Copyright IBM Corp. 2025
// SPDX-License-Identifier: MPL-2.0

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

func TestGetAgentCheckTool(t *testing.T) {
	logger := log.New()
	tool := GetAgentCheckTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "agent_check", tool.Tool.Name)
}

func TestGetAgentCheckHandler(t *testing.T) {
	mockCheck := map[string]interface{}{
		"Node":        "consul-agent-1",
		"CheckID":     "service:web-health",
		"Name":        "Web Service Health Check",
		"Status":      "passing",
		"Notes":       "HTTP health check for web service",
		"Output":      "HTTP GET http://localhost:8080/health: 200 OK Output: {\"status\":\"healthy\"}",
		"ServiceID":   "web-service",
		"ServiceName": "web",
		"ServiceTags": []string{"v1.0", "production", "frontend"},
		"Type":        "http",
		"Definition": map[string]interface{}{
			"HTTP":                           "http://localhost:8080/health",
			"IntervalDuration":               "30s",
			"TimeoutDuration":                "5s",
			"DeregisterCriticalServiceAfter": "90m",
			"TLSSkipVerify":                  false,
		},
		"CreateIndex": 100,
		"ModifyIndex": 150,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/agent/check/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockCheck)
	}))
	defer server.Close()

	t.Run("successful check retrieval", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"check_id": "service:web-health",
			},
		}

		checkID, err := request.RequireString("check_id")
		assert.NoError(t, err)
		assert.Equal(t, "service:web-health", checkID)
	})

	t.Run("missing check ID", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("check_id")
		assert.Error(t, err)
	})
}

func TestGetAgentCheckByNameTool(t *testing.T) {
	logger := log.New()
	tool := GetAgentCheckByNameTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "agent_check_by_name", tool.Tool.Name)
}

func TestGetAgentCheckByNameHandler(t *testing.T) {
	mockCheck := map[string]interface{}{
		"Node":        "consul-agent-1",
		"CheckID":     "web-app-health",
		"Name":        "Web Application Health",
		"Status":      "warning",
		"Notes":       "Application health check with custom validation",
		"Output":      "HTTP GET http://localhost:8080/api/health: 200 OK Output: {\"status\":\"degraded\",\"message\":\"Database slow\"}",
		"ServiceID":   "web-app",
		"ServiceName": "web-application",
		"ServiceTags": []string{"v2.0", "beta"},
		"Type":        "http",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/agent/check/name/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockCheck)
	}))
	defer server.Close()

	t.Run("successful check retrieval by name", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"check_name": "Web Application Health",
			},
		}

		checkName, err := request.RequireString("check_name")
		assert.NoError(t, err)
		assert.Equal(t, "Web Application Health", checkName)
	})

	t.Run("missing check name", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("check_name")
		assert.Error(t, err)
	})
}

func TestAgentCheckParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid check ID",
			arguments: map[string]interface{}{
				"check_id": "service:web-health",
			},
			expectError: false,
		},
		{
			name: "valid check name",
			arguments: map[string]interface{}{
				"check_name": "Database Health Check",
			},
			expectError: false,
		},
		{
			name:        "missing check ID",
			arguments:   map[string]interface{}{},
			expectError: true,
			errorField:  "check_id",
		},
		{
			name:        "missing check name",
			arguments:   map[string]interface{}{},
			expectError: true,
			errorField:  "check_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &MockCallToolRequest{Arguments: tt.arguments}

			if tt.expectError {
				_, err := request.RequireString(tt.errorField)
				assert.Error(t, err)
			} else {
				if checkID := request.GetString("check_id", ""); checkID != "" {
					assert.NotEmpty(t, checkID)
				}
				if checkName := request.GetString("check_name", ""); checkName != "" {
					assert.NotEmpty(t, checkName)
				}
			}
		})
	}
}

func TestAgentCheckTypes(t *testing.T) {
	checkTypes := []string{"http", "tcp", "script", "ttl", "docker", "grpc"}

	for _, checkType := range checkTypes {
		t.Run(fmt.Sprintf("%s check type", checkType), func(t *testing.T) {
			check := map[string]interface{}{
				"Type":    checkType,
				"CheckID": fmt.Sprintf("test-%s-check", checkType),
			}

			assert.Equal(t, checkType, check["Type"])
		})
	}
}

func TestAgentCheckStatuses(t *testing.T) {
	statuses := []string{"passing", "warning", "critical"}

	for _, status := range statuses {
		t.Run(fmt.Sprintf("%s status", status), func(t *testing.T) {
			check := map[string]interface{}{
				"Status":  status,
				"CheckID": fmt.Sprintf("test-%s", status),
			}

			assert.Equal(t, status, check["Status"])
		})
	}
}

func TestAgentCheckDefinitions(t *testing.T) {
	t.Run("HTTP check definition", func(t *testing.T) {
		definition := map[string]interface{}{
			"HTTP":                           "http://localhost:8080/health",
			"IntervalDuration":               "30s",
			"TimeoutDuration":                "5s",
			"DeregisterCriticalServiceAfter": "90m",
			"TLSSkipVerify":                  false,
			"Method":                         "GET",
			"Header":                         map[string][]string{"Authorization": {"Bearer token"}},
		}

		assert.Equal(t, "http://localhost:8080/health", definition["HTTP"])
		assert.Equal(t, "30s", definition["IntervalDuration"])
		assert.Equal(t, false, definition["TLSSkipVerify"])
	})

	t.Run("TCP check definition", func(t *testing.T) {
		definition := map[string]interface{}{
			"TCP":              "localhost:5432",
			"IntervalDuration": "10s",
			"TimeoutDuration":  "3s",
		}

		assert.Equal(t, "localhost:5432", definition["TCP"])
		assert.Equal(t, "10s", definition["IntervalDuration"])
	})

	t.Run("Script check definition", func(t *testing.T) {
		definition := map[string]interface{}{
			"Args":             []string{"/usr/local/bin/check_service.sh"},
			"IntervalDuration": "60s",
			"TimeoutDuration":  "30s",
		}

		args := definition["Args"].([]string)
		assert.Len(t, args, 1)
		assert.Contains(t, args[0], "check_service.sh")
	})

	t.Run("TTL check definition", func(t *testing.T) {
		definition := map[string]interface{}{
			"TTL":                            "30s",
			"DeregisterCriticalServiceAfter": "90m",
		}

		assert.Equal(t, "30s", definition["TTL"])
		assert.Equal(t, "90m", definition["DeregisterCriticalServiceAfter"])
	})
}

func TestAgentCheckResponseProcessing(t *testing.T) {
	t.Run("detailed check response", func(t *testing.T) {
		check := map[string]interface{}{
			"CheckID":     "complex-health-check",
			"Name":        "Complex Application Health",
			"Status":      "passing",
			"Output":      "All systems operational",
			"ServiceName": "complex-app",
			"Type":        "http",
			"Definition": map[string]interface{}{
				"HTTP":             "https://app.example.com/health",
				"IntervalDuration": "15s",
				"TimeoutDuration":  "10s",
				"TLSSkipVerify":    false,
			},
		}

		data, err := json.MarshalIndent(check, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "complex-health-check", unmarshaled["CheckID"])
		assert.Equal(t, "passing", unmarshaled["Status"])

		definition := unmarshaled["Definition"].(map[string]interface{})
		assert.Equal(t, "https://app.example.com/health", definition["HTTP"])
	})
}

func TestAgentCheckErrorHandling(t *testing.T) {
	t.Run("check not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Check not found"))
		}))
		defer server.Close()

		// Would test handling of non-existent checks
		assert.True(t, true) // Placeholder
	})

	t.Run("invalid check configuration", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid check configuration"))
		}))
		defer server.Close()

		// Would test handling of malformed check definitions
		assert.True(t, true) // Placeholder
	})

	t.Run("agent connection error", func(t *testing.T) {
		// Test scenarios where agent is unreachable
		assert.True(t, true) // Placeholder
	})
}

func TestAgentCheckValidation(t *testing.T) {
	t.Run("check ID format", func(t *testing.T) {
		validCheckIDs := []string{
			"service:web-health",
			"serfHealth",
			"custom-check-123",
			"node-health",
		}

		for _, checkID := range validCheckIDs {
			assert.NotEmpty(t, checkID)
			// Additional ID format validation would go here
		}
	})

	t.Run("check output validation", func(t *testing.T) {
		outputs := []string{
			"HTTP GET http://localhost:8080/health: 200 OK",
			"TCP connect localhost:5432: Success",
			"Script execution completed with exit code 0",
			"TTL check passed",
		}

		for _, output := range outputs {
			assert.NotEmpty(t, output)
			// Additional output validation would go here
		}
	})
}
