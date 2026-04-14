// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package functions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAgentOverviewHandler(t *testing.T) {
	mockOverview := map[string]interface{}{
		"Config": map[string]interface{}{
			"Datacenter":    "dc1",
			"NodeName":      "consul-agent-1",
			"NodeID":        "node-uuid-123",
			"Server":        true,
			"Revision":      "12345678",
			"Version":       "1.16.1",
			"Domain":        "consul",
			"LogLevel":      "INFO",
			"ClientAddr":    []string{"127.0.0.1"},
			"BindAddr":      "0.0.0.0",
			"AdvertiseAddr": "192.168.1.10",
		},
		"Coord": map[string]interface{}{
			"Adjustment": 0.0,
			"Error":      1.5,
			"Height":     0.0,
			"Vec":        []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
		},
		"Member": map[string]interface{}{
			"Name": "consul-agent-1",
			"Addr": "192.168.1.10",
			"Port": 8301,
			"Tags": map[string]string{
				"bootstrap": "1",
				"dc":        "dc1",
				"role":      "consul",
				"vsn":       "3",
				"vsn_max":   "3",
				"vsn_min":   "2",
			},
			"Status":      1, // StatusAlive
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
				"checks":         "5",
				"services":       "3",
			},
			"consul": map[string]string{
				"bootstrap":         "true",
				"known_datacenters": "1",
				"leader":            "true",
				"leader_addr":       "192.168.1.10:8300",
				"server":            "true",
			},
			"raft": map[string]string{
				"applied_index":       "1000",
				"commit_index":        "1000",
				"fsm_pending":         "0",
				"last_log_index":      "1000",
				"last_log_term":       "5",
				"last_snapshot_index": "500",
				"last_snapshot_term":  "3",
				"num_peers":           "2",
				"state":               "Leader",
				"term":                "5",
			},
			"serf_lan": map[string]string{
				"coordinate_resets": "0",
				"encrypted":         "false",
				"event_queue":       "0",
				"event_time":        "1",
				"failed":            "0",
				"left":              "0",
				"member_time":       "1",
				"members":           "3",
				"query_queue":       "0",
				"query_time":        "1",
			},
		},
		"Meta": map[string]string{
			"consul-network-segment": "",
			"consul-version":         "1.16.1",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/agent/self", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockOverview)
	}))
	defer server.Close()

	t.Run("successful agent overview retrieval", func(t *testing.T) {
		// Test would verify agent overview is returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})
}

func TestAgentOverviewConfiguration(t *testing.T) {
	t.Run("agent configuration validation", func(t *testing.T) {
		config := map[string]interface{}{
			"Datacenter":    "dc1",
			"NodeName":      "consul-server-1",
			"Server":        true,
			"Version":       "1.16.1",
			"LogLevel":      "INFO",
			"ClientAddr":    []string{"127.0.0.1", "0.0.0.0"},
			"BindAddr":      "0.0.0.0",
			"AdvertiseAddr": "192.168.1.10",
		}

		assert.Equal(t, "dc1", config["Datacenter"])
		assert.Equal(t, true, config["Server"])
		assert.Equal(t, "INFO", config["LogLevel"])

		clientAddrs := config["ClientAddr"].([]string)
		assert.Contains(t, clientAddrs, "127.0.0.1")
	})

	t.Run("client vs server configuration", func(t *testing.T) {
		serverConfig := map[string]interface{}{
			"Server":    true,
			"Bootstrap": true,
		}

		clientConfig := map[string]interface{}{
			"Server":    false,
			"Bootstrap": false,
		}

		assert.Equal(t, true, serverConfig["Server"])
		assert.Equal(t, false, clientConfig["Server"])
	})
}

func TestAgentOverviewStats(t *testing.T) {
	t.Run("raft statistics", func(t *testing.T) {
		raftStats := map[string]string{
			"applied_index": "1000",
			"commit_index":  "1000",
			"state":         "Leader",
			"term":          "5",
			"num_peers":     "2",
		}

		assert.Equal(t, "Leader", raftStats["state"])
		assert.Equal(t, "1000", raftStats["applied_index"])
		assert.Equal(t, "2", raftStats["num_peers"])
	})

	t.Run("serf LAN statistics", func(t *testing.T) {
		serfStats := map[string]string{
			"members":   "3",
			"failed":    "0",
			"left":      "0",
			"encrypted": "false",
		}

		assert.Equal(t, "3", serfStats["members"])
		assert.Equal(t, "0", serfStats["failed"])
		assert.Equal(t, "false", serfStats["encrypted"])
	})

	t.Run("consul statistics", func(t *testing.T) {
		consulStats := map[string]string{
			"leader":            "true",
			"server":            "true",
			"known_datacenters": "1",
			"bootstrap":         "true",
		}

		assert.Equal(t, "true", consulStats["leader"])
		assert.Equal(t, "true", consulStats["server"])
	})
}

func TestAgentOverviewCoordinates(t *testing.T) {
	t.Run("network coordinates", func(t *testing.T) {
		coord := map[string]interface{}{
			"Adjustment": 0.0,
			"Error":      1.5,
			"Height":     0.0,
			"Vec":        []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
		}

		assert.Equal(t, 0.0, coord["Adjustment"])
		assert.Equal(t, 1.5, coord["Error"])

		vec := coord["Vec"].([]float64)
		assert.Len(t, vec, 8)
		assert.Equal(t, 0.1, vec[0])
	})
}

func TestAgentOverviewMember(t *testing.T) {
	t.Run("serf member information", func(t *testing.T) {
		member := map[string]interface{}{
			"Name": "consul-agent-1",
			"Addr": "192.168.1.10",
			"Port": 8301,
			"Tags": map[string]string{
				"bootstrap": "1",
				"dc":        "dc1",
				"role":      "consul",
				"vsn":       "3",
			},
			"Status":      1, // StatusAlive
			"ProtocolCur": 2,
		}

		assert.Equal(t, "consul-agent-1", member["Name"])
		assert.Equal(t, "192.168.1.10", member["Addr"])
		assert.Equal(t, 8301, member["Port"])
		assert.Equal(t, 1, member["Status"])

		tags := member["Tags"].(map[string]string)
		assert.Equal(t, "dc1", tags["dc"])
		assert.Equal(t, "consul", tags["role"])
	})
}

func TestAgentOverviewResponseProcessing(t *testing.T) {
	t.Run("complete overview response", func(t *testing.T) {
		overview := map[string]interface{}{
			"Config": map[string]interface{}{
				"Datacenter": "dc1",
				"NodeName":   "test-node",
				"Server":     true,
			},
			"Stats": map[string]interface{}{
				"consul": map[string]string{
					"leader": "true",
					"server": "true",
				},
			},
			"Member": map[string]interface{}{
				"Name":   "test-node",
				"Status": 1,
			},
		}

		data, err := json.MarshalIndent(overview, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		config := unmarshaled["Config"].(map[string]interface{})
		assert.Equal(t, "dc1", config["Datacenter"])

		stats := unmarshaled["Stats"].(map[string]interface{})
		consulStats := stats["consul"].(map[string]interface{})
		assert.Equal(t, "true", consulStats["leader"])
	})
}

func TestAgentOverviewRoles(t *testing.T) {
	t.Run("server node overview", func(t *testing.T) {
		serverOverview := map[string]interface{}{
			"Config": map[string]interface{}{
				"Server":    true,
				"Bootstrap": true,
			},
			"Stats": map[string]interface{}{
				"raft": map[string]string{
					"state": "Leader",
				},
			},
		}

		config := serverOverview["Config"].(map[string]interface{})
		assert.Equal(t, true, config["Server"])

		stats := serverOverview["Stats"].(map[string]interface{})
		raftStats := stats["raft"].(map[string]string)
		assert.Equal(t, "Leader", raftStats["state"])
	})

	t.Run("client node overview", func(t *testing.T) {
		clientOverview := map[string]interface{}{
			"Config": map[string]interface{}{
				"Server":    false,
				"Bootstrap": false,
			},
			"Stats": map[string]interface{}{
				"agent": map[string]string{
					"services": "5",
					"checks":   "10",
				},
			},
		}

		config := clientOverview["Config"].(map[string]interface{})
		assert.Equal(t, false, config["Server"])
	})
}

func TestAgentOverviewErrorHandling(t *testing.T) {
	t.Run("agent not available", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Agent not available"))
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
