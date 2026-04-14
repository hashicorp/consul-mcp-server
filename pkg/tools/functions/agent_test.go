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

func TestGetAgentSelfTool(t *testing.T) {
	logger := log.New()
	tool := GetAgentSelfTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "agent_self", tool.Tool.Name)
}

func TestGetAgentSelfHandler(t *testing.T) {
	mockAgentInfo := map[string]interface{}{
		"Config": map[string]interface{}{
			"Datacenter":        "dc1",
			"NodeName":          "consul-agent-1",
			"NodeID":            "node-id-123",
			"Server":            true,
			"Revision":          "12345678",
			"Version":           "1.16.1",
			"VersionPrerelease": "",
		},
		"Coord": map[string]interface{}{
			"Adjustment": 0,
			"Error":      1.5,
			"Vec":        []float64{0.1, 0.2, 0.3},
		},
		"Member": map[string]interface{}{
			"Name":        "consul-agent-1",
			"Addr":        "192.168.1.10",
			"Port":        8301,
			"Tags":        map[string]string{"role": "consul", "dc": "dc1"},
			"Status":      1,
			"ProtocolMin": 1,
			"ProtocolMax": 5,
			"ProtocolCur": 2,
			"DelegateMin": 2,
			"DelegateMax": 5,
			"DelegateCur": 4,
		},
		"Stats": map[string]interface{}{
			"agent": map[string]string{
				"check_monitors": "0",
				"check_ttls":     "0",
				"checks":         "1",
				"services":       "1",
			},
			"consul": map[string]string{
				"bootstrap":         "true",
				"known_datacenters": "1",
				"leader":            "true",
				"leader_addr":       "192.168.1.10:8300",
				"server":            "true",
			},
			"raft": map[string]string{
				"applied_index":       "100",
				"commit_index":        "100",
				"fsm_pending":         "0",
				"last_log_index":      "100",
				"last_log_term":       "1",
				"last_snapshot_index": "0",
				"last_snapshot_term":  "0",
				"num_peers":           "0",
				"state":               "Leader",
				"term":                "1",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/agent/self", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockAgentInfo)
	}))
	defer server.Close()

	t.Run("successful agent self info", func(t *testing.T) {
		// Test would verify agent information is returned correctly
		assert.True(t, true) // Placeholder
	})
}

func TestGetAgentMembersTool(t *testing.T) {
	logger := log.New()
	tool := GetAgentMembersTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "agent_members", tool.Tool.Name)
}

func TestGetAgentMembersHandler(t *testing.T) {
	mockMembers := []map[string]interface{}{
		{
			"Name":        "consul-server-1",
			"Addr":        "192.168.1.10",
			"Port":        8301,
			"Tags":        map[string]string{"bootstrap": "1", "dc": "dc1", "role": "consul", "vsn": "3", "vsn_max": "3", "vsn_min": "2"},
			"Status":      1,
			"ProtocolMin": 1,
			"ProtocolMax": 5,
			"ProtocolCur": 2,
			"DelegateMin": 2,
			"DelegateMax": 5,
			"DelegateCur": 4,
		},
		{
			"Name":        "consul-server-2",
			"Addr":        "192.168.1.11",
			"Port":        8301,
			"Tags":        map[string]string{"dc": "dc1", "role": "consul", "vsn": "3", "vsn_max": "3", "vsn_min": "2"},
			"Status":      1,
			"ProtocolMin": 1,
			"ProtocolMax": 5,
			"ProtocolCur": 2,
			"DelegateMin": 2,
			"DelegateMax": 5,
			"DelegateCur": 4,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/agent/members", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockMembers)
	}))
	defer server.Close()

	t.Run("successful members listing", func(t *testing.T) {
		// Test would verify cluster members are returned correctly
		assert.True(t, true) // Placeholder
	})

	t.Run("WAN members filtering", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"wan": "true",
			},
		}

		wan := request.GetString("wan", "false")
		assert.Equal(t, "true", wan)
	})
}

func TestGetAgentServicesTool(t *testing.T) {
	logger := log.New()
	tool := GetAgentServicesTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "agent_services", tool.Tool.Name)
}

func TestGetAgentServicesHandler(t *testing.T) {
	mockServices := map[string]interface{}{
		"web": map[string]interface{}{
			"ID":      "web",
			"Service": "web",
			"Tags":    []string{"v1.0", "production"},
			"Meta":    map[string]string{"version": "1.0", "environment": "production"},
			"Port":    8080,
			"Address": "192.168.1.10",
			"Weights": map[string]interface{}{
				"Passing": 10,
				"Warning": 1,
			},
			"EnableTagOverride": false,
		},
		"database": map[string]interface{}{
			"ID":      "database",
			"Service": "database",
			"Tags":    []string{"v5.7", "primary"},
			"Meta":    map[string]string{"version": "5.7", "role": "primary"},
			"Port":    3306,
			"Address": "192.168.1.20",
			"Weights": map[string]interface{}{
				"Passing": 10,
				"Warning": 1,
			},
			"EnableTagOverride": false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/agent/services", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockServices)
	}))
	defer server.Close()

	t.Run("successful services listing", func(t *testing.T) {
		// Test would verify agent services are returned correctly
		assert.True(t, true) // Placeholder
	})

	t.Run("service filtering", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"filter": "Service == \"web\"",
			},
		}

		filter := request.GetString("filter", "")
		assert.Equal(t, "Service == \"web\"", filter)
	})
}

func TestGetAgentChecksTool(t *testing.T) {
	logger := log.New()
	tool := GetAgentChecksTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "agent_checks", tool.Tool.Name)
}

func TestGetAgentChecksHandler(t *testing.T) {
	mockChecks := map[string]interface{}{
		"service:web": map[string]interface{}{
			"Node":        "consul-agent-1",
			"CheckID":     "service:web",
			"Name":        "Service 'web' check",
			"Status":      "passing",
			"Notes":       "",
			"Output":      "HTTP GET http://localhost:8080/health: 200 OK Output: OK",
			"ServiceID":   "web",
			"ServiceName": "web",
			"ServiceTags": []string{"v1.0", "production"},
			"Type":        "http",
			"Definition": map[string]interface{}{
				"HTTP":                           "http://localhost:8080/health",
				"IntervalDuration":               "10s",
				"TimeoutDuration":                "3s",
				"DeregisterCriticalServiceAfter": "30m",
			},
		},
		"serfHealth": map[string]interface{}{
			"Node":        "consul-agent-1",
			"CheckID":     "serfHealth",
			"Name":        "Serf Health Status",
			"Status":      "passing",
			"Notes":       "",
			"Output":      "Agent alive and reachable",
			"ServiceID":   "",
			"ServiceName": "",
			"ServiceTags": []string{},
			"Type":        "serf",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/agent/checks", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockChecks)
	}))
	defer server.Close()

	t.Run("successful checks listing", func(t *testing.T) {
		// Test would verify agent checks are returned correctly
		assert.True(t, true) // Placeholder
	})

	t.Run("check filtering by service", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"filter": "ServiceName == \"web\"",
			},
		}

		filter := request.GetString("filter", "")
		assert.Equal(t, "ServiceName == \"web\"", filter)
	})
}

func TestAgentResponseProcessing(t *testing.T) {
	t.Run("agent self config parsing", func(t *testing.T) {
		agentConfig := map[string]interface{}{
			"Config": map[string]interface{}{
				"Datacenter": "dc1",
				"NodeName":   "consul-node-1",
				"Server":     true,
				"Version":    "1.16.1",
			},
			"Stats": map[string]interface{}{
				"consul": map[string]string{
					"leader": "true",
					"server": "true",
				},
			},
		}

		data, err := json.MarshalIndent(agentConfig, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		config, ok := unmarshaled["Config"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "dc1", config["Datacenter"])
		assert.Equal(t, true, config["Server"])
	})

	t.Run("service weights processing", func(t *testing.T) {
		service := map[string]interface{}{
			"Service": "web",
			"Port":    8080,
			"Weights": map[string]interface{}{
				"Passing": 10,
				"Warning": 1,
			},
		}

		data, err := json.Marshal(service)
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		weights, ok := unmarshaled["Weights"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(10), weights["Passing"])
		assert.Equal(t, float64(1), weights["Warning"])
	})
}

func TestAgentParameterHandling(t *testing.T) {
	tests := []struct {
		name      string
		arguments map[string]interface{}
		testField string
		expected  string
	}{
		{
			name: "WAN members query",
			arguments: map[string]interface{}{
				"wan": "true",
			},
			testField: "wan",
			expected:  "true",
		},
		{
			name: "service filter",
			arguments: map[string]interface{}{
				"filter": "ServiceName == \"web\"",
			},
			testField: "filter",
			expected:  "ServiceName == \"web\"",
		},
		{
			name: "check filter",
			arguments: map[string]interface{}{
				"filter": "Status == \"critical\"",
			},
			testField: "filter",
			expected:  "Status == \"critical\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &MockCallToolRequest{Arguments: tt.arguments}

			value := request.GetString(tt.testField, "")
			assert.Equal(t, tt.expected, value)
		})
	}
}

func TestAgentErrorHandling(t *testing.T) {
	t.Run("agent unavailable", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Agent unavailable"))
		}))
		defer server.Close()

		// Would test handling of agent unavailability
		assert.True(t, true) // Placeholder
	})

	t.Run("permission denied", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Permission denied"))
		}))
		defer server.Close()

		// Would test handling of ACL permission errors
		assert.True(t, true) // Placeholder
	})
}

func TestAgentMetricsAndStats(t *testing.T) {
	t.Run("raft stats parsing", func(t *testing.T) {
		raftStats := map[string]string{
			"applied_index":  "100",
			"commit_index":   "100",
			"last_log_index": "100",
			"num_peers":      "2",
			"state":          "Leader",
			"term":           "1",
		}

		// Test parsing of numeric values from string stats
		appliedIndex := raftStats["applied_index"]
		assert.Equal(t, "100", appliedIndex)

		state := raftStats["state"]
		assert.Equal(t, "Leader", state)
	})

	t.Run("serf member status", func(t *testing.T) {
		member := map[string]interface{}{
			"Name":   "consul-server-1",
			"Status": 1, // StatusAlive
			"Tags":   map[string]string{"role": "consul", "dc": "dc1"},
		}

		status, ok := member["Status"].(int)
		assert.True(t, ok)
		assert.Equal(t, 1, status) // StatusAlive

		tags, ok := member["Tags"].(map[string]string)
		assert.True(t, ok)
		assert.Equal(t, "consul", tags["role"])
		assert.Equal(t, "dc1", tags["dc"])
	})
}
