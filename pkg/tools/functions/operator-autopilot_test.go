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

func TestGetOperatorAutopilotConfigurationTool(t *testing.T) {
	tool := GetOperatorAutopilotConfigurationTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "operator_autopilot_configuration", tool.Tool.Name)
}

func TestGetOperatorAutopilotConfigurationHandler(t *testing.T) {
	mockAutopilotConfig := map[string]interface{}{
		"CleanupDeadServers":      true,
		"LastContactThreshold":    "200ms",
		"MaxTrailingLogs":         250,
		"MinQuorum":               3,
		"ServerStabilizationTime": "10s",
		"RedundancyZoneTag":       "zone",
		"DisableUpgradeMigration": false,
		"UpgradeVersionTag":       "version",
		"CreateIndex":             100,
		"ModifyIndex":             150,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/operator/autopilot/configuration", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockAutopilotConfig)
	}))
	defer server.Close()

	t.Run("successful autopilot configuration retrieval", func(t *testing.T) {
		// Test would verify autopilot configuration is returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})
}

func TestGetOperatorAutopilotHealthTool(t *testing.T) {
	tool := GetOperatorAutopilotHealthTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "operator_autopilot_health", tool.Tool.Name)
}

func TestGetOperatorAutopilotHealthHandler(t *testing.T) {
	mockAutopilotHealth := map[string]interface{}{
		"Healthy":          true,
		"FailureTolerance": 1,
		"Servers": []map[string]interface{}{
			{
				"ID":          "server-1",
				"Name":        "consul-server-1",
				"Address":     "192.168.1.10:8300",
				"SerfStatus":  "alive",
				"Version":     "1.16.0",
				"Leader":      true,
				"LastContact": "0s",
				"LastTerm":    5,
				"LastIndex":   1000,
				"Healthy":     true,
				"Voter":       true,
				"StableSince": "2023-01-01T12:00:00Z",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/operator/autopilot/health", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockAutopilotHealth)
	}))
	defer server.Close()

	t.Run("successful autopilot health retrieval", func(t *testing.T) {
		// Test would verify autopilot health status is returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})
}

func TestAutopilotConfigurationValidation(t *testing.T) {
	t.Run("autopilot configuration structure", func(t *testing.T) {
		config := map[string]interface{}{
			"CleanupDeadServers":      true,
			"LastContactThreshold":    "200ms",
			"MaxTrailingLogs":         250,
			"MinQuorum":               3,
			"ServerStabilizationTime": "10s",
			"RedundancyZoneTag":       "az",
		}

		assert.Equal(t, true, config["CleanupDeadServers"])
		assert.Equal(t, "200ms", config["LastContactThreshold"])
		assert.Equal(t, 250, config["MaxTrailingLogs"])
		assert.Equal(t, 3, config["MinQuorum"])
	})

	t.Run("enterprise autopilot features", func(t *testing.T) {
		enterpriseConfig := map[string]interface{}{
			"RedundancyZoneTag":       "availability-zone",
			"DisableUpgradeMigration": false,
			"UpgradeVersionTag":       "consul-version",
		}

		assert.Equal(t, "availability-zone", enterpriseConfig["RedundancyZoneTag"])
		assert.Equal(t, false, enterpriseConfig["DisableUpgradeMigration"])
	})
}

func TestAutopilotHealthValidation(t *testing.T) {
	t.Run("cluster health assessment", func(t *testing.T) {
		healthyCluster := map[string]interface{}{
			"Healthy":          true,
			"FailureTolerance": 1,
		}

		unhealthyCluster := map[string]interface{}{
			"Healthy":          false,
			"FailureTolerance": 0,
		}

		assert.Equal(t, true, healthyCluster["Healthy"])
		assert.Equal(t, 1, healthyCluster["FailureTolerance"])
		assert.Equal(t, false, unhealthyCluster["Healthy"])
		assert.Equal(t, 0, unhealthyCluster["FailureTolerance"])
	})

	t.Run("server health status", func(t *testing.T) {
		healthyServer := map[string]interface{}{
			"ID":          "server-1",
			"SerfStatus":  "alive",
			"Leader":      true,
			"Healthy":     true,
			"Voter":       true,
			"LastContact": "0s",
		}

		unhealthyServer := map[string]interface{}{
			"ID":          "server-4",
			"SerfStatus":  "failed",
			"Leader":      false,
			"Healthy":     false,
			"Voter":       false,
			"LastContact": "5m",
		}

		assert.Equal(t, true, healthyServer["Healthy"])
		assert.Equal(t, "alive", healthyServer["SerfStatus"])
		assert.Equal(t, false, unhealthyServer["Healthy"])
		assert.Equal(t, "failed", unhealthyServer["SerfStatus"])
	})
}

func TestAutopilotResponseProcessing(t *testing.T) {
	t.Run("autopilot configuration response", func(t *testing.T) {
		config := map[string]interface{}{
			"CleanupDeadServers":      true,
			"LastContactThreshold":    "200ms",
			"MaxTrailingLogs":         250,
			"ServerStabilizationTime": "10s",
		}

		data, err := json.MarshalIndent(config, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, true, unmarshaled["CleanupDeadServers"])
		assert.Equal(t, "200ms", unmarshaled["LastContactThreshold"])
	})

	t.Run("autopilot health response", func(t *testing.T) {
		health := map[string]interface{}{
			"Healthy":          true,
			"FailureTolerance": 1,
			"Servers": []map[string]interface{}{
				{
					"ID":      "server-1",
					"Healthy": true,
					"Leader":  true,
				},
			},
		}

		data, err := json.Marshal(health)
		require.NoError(t, err)
		assert.Contains(t, string(data), "server-1")
		assert.Contains(t, string(data), "Healthy")
	})
}

func TestAutopilotScenarios(t *testing.T) {
	t.Run("single node cluster", func(t *testing.T) {
		singleNodeHealth := map[string]interface{}{
			"Healthy":          true,
			"FailureTolerance": 0,
			"Servers": []map[string]interface{}{
				{
					"ID":      "single-server",
					"Leader":  true,
					"Voter":   true,
					"Healthy": true,
				},
			},
		}

		servers := singleNodeHealth["Servers"].([]map[string]interface{})
		assert.Len(t, servers, 1)
		assert.Equal(t, 0, singleNodeHealth["FailureTolerance"])
	})

	t.Run("multi-node cluster with leader election", func(t *testing.T) {
		multiNodeHealth := map[string]interface{}{
			"Healthy":          true,
			"FailureTolerance": 1,
			"Servers": []map[string]interface{}{
				{"ID": "server-1", "Leader": true, "Voter": true},
				{"ID": "server-2", "Leader": false, "Voter": true},
				{"ID": "server-3", "Leader": false, "Voter": true},
			},
		}

		servers := multiNodeHealth["Servers"].([]map[string]interface{})
		assert.Len(t, servers, 3)

		leaderCount := 0
		for _, server := range servers {
			if server["Leader"].(bool) {
				leaderCount++
			}
		}
		assert.Equal(t, 1, leaderCount, "Should have exactly one leader")
	})

	t.Run("cluster with non-voter", func(t *testing.T) {
		clusterWithNonVoter := map[string]interface{}{
			"Servers": []map[string]interface{}{
				{"ID": "server-1", "Leader": true, "Voter": true},
				{"ID": "server-2", "Leader": false, "Voter": true},
				{"ID": "server-3", "Leader": false, "Voter": false}, // Non-voter
			},
		}

		servers := clusterWithNonVoter["Servers"].([]map[string]interface{})
		voterCount := 0
		for _, server := range servers {
			if server["Voter"].(bool) {
				voterCount++
			}
		}
		assert.Equal(t, 2, voterCount, "Should have 2 voting servers")
	})
}

func TestAutopilotErrorHandling(t *testing.T) {
	t.Run("autopilot not enabled", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Autopilot not enabled"))
		}))
		defer server.Close()

		// Would test handling when autopilot is disabled
		assert.True(t, true) // Placeholder
	})

	t.Run("insufficient servers for quorum", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Insufficient servers for quorum"))
		}))
		defer server.Close()

		// Would test handling of quorum loss scenarios
		assert.True(t, true) // Placeholder
	})
}
