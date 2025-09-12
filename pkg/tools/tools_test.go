// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"testing"

	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRegisterTools(t *testing.T) {
	// Create a mock MCP server
	mcpServer := &mockMCPServer{
		tools: make(map[string]server.ServerTool),
	}

	// Test tool registration
	RegisterTools(mcpServer, log.New())

	// Verify that tools were registered
	assert.Greater(t, len(mcpServer.tools), 0, "Expected tools to be registered")

	// Check for some expected tools
	expectedTools := []string{
		"catalog_services",
		"catalog_nodes",
		"health_service",
		"kv_get",
		"acl_tokens",
	}

	for _, toolName := range expectedTools {
		_, exists := mcpServer.tools[toolName]
		assert.True(t, exists, "Expected tool %s to be registered", toolName)
	}

	// Check the number of tools registered
	assert.Equal(t, len(mcpServer.tools), 85)
}

// Mock MCP Server for testing
type mockMCPServer struct {
	tools map[string]server.ServerTool
}

func (m *mockMCPServer) AddTool(tool mcp.Tool, handler server.ToolHandlerFunc) {
	m.tools[tool.Name] = server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}
