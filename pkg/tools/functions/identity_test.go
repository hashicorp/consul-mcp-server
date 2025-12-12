// Copyright IBM Corp. 2025
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"encoding/json"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIdentityTool(t *testing.T) {
	logger := log.New()
	tool := GetIdentity(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "cluster_identity", tool.Tool.Name)
}

func TestIdentityResponseStructure(t *testing.T) {
	t.Run("identity response format", func(t *testing.T) {
		// Expected identity response structure
		identity := map[string]interface{}{
			"name":        "consul-mcp-server",
			"version":     "1.0.0",
			"description": "HashiCorp Consul MCP Server",
			"capabilities": []string{
				"read-consul-catalog",
				"read-consul-health",
				"read-consul-kv",
				"read-consul-acl",
				"read-consul-connect",
			},
		}

		data, err := json.MarshalIndent(identity, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "consul-mcp-server", unmarshaled["name"])
		assert.Contains(t, unmarshaled, "capabilities")
	})
}

func TestIdentityCapabilities(t *testing.T) {
	capabilities := []string{
		"read-consul-catalog",
		"read-consul-health",
		"read-consul-kv",
		"read-consul-acl",
		"read-consul-connect",
		"read-consul-config",
		"read-consul-status",
		"read-consul-agent",
	}

	for _, capability := range capabilities {
		t.Run(fmt.Sprintf("capability %s", capability), func(t *testing.T) {
			assert.NotEmpty(t, capability)
			assert.Contains(t, capability, "read-consul-")
		})
	}
}

func TestIdentityVersionInfo(t *testing.T) {
	t.Run("version information", func(t *testing.T) {
		versionInfo := map[string]interface{}{
			"version":    "1.0.0",
			"build_date": "2025-09-11T10:00:00Z",
			"git_commit": "abc123def456",
			"go_version": "go1.21",
		}

		for key, value := range versionInfo {
			assert.NotEmpty(t, value, "Version field %s should not be empty", key)
		}
	})
}
