// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package e2e

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"testing"
	"time"

	mcpClient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	buildDockerImage(t)

	// Ensure all test containers are cleaned up at the end
	t.Cleanup(func() {
		cleanupAllTestContainers(t)
	})

	testCases := []struct {
		name          string
		clientFactory func(t *testing.T) (mcpClient.MCPClient, func())
	}{
		{"Stdio", createStdioClient},
		{"HTTP", createHTTPClient},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, cleanup := tc.clientFactory(t)
			defer cleanup()
			runTestSuite(t, client, tc.name)
		})
	}
}

// ensureClientInitialized ensures the MCP client is initialized before running tool tests
func ensureClientInitialized(t *testing.T, client mcpClient.MCPClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := mcp.InitializeRequest{}
	request.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	request.Params.ClientInfo = mcp.Implementation{
		Name:    "e2e-test-client",
		Version: "0.0.1",
	}

	result, err := client.Initialize(ctx, request)
	if err != nil {
		t.Fatalf("Failed to initialize MCP client: %v", err)
	}
	t.Logf("Initialized with server: %s %s", result.ServerInfo.Name, result.ServerInfo.Version)
	require.Equal(t, "consul-mcp-server", result.ServerInfo.Name)
}

// runTestSuite executes all test cases against the provided client
func runTestSuite(t *testing.T, client mcpClient.MCPClient, transportName string) {
	t.Run("Initialize", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		request := mcp.InitializeRequest{}
		request.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		request.Params.ClientInfo = mcp.Implementation{
			Name:    "e2e-test-client",
			Version: "0.0.1",
		}

		result, err := client.Initialize(ctx, request)
		if err != nil {
			log.Fatalf("Failed to initialize: %v", err)
		}
		fmt.Printf(
			"Initialized with server: %s %s\n\n",
			result.ServerInfo.Name,
			result.ServerInfo.Version,
		)
		require.Equal(t, "consul-mcp-server", result.ServerInfo.Name)
	})

	t.Run("Tool count Tests", func(t *testing.T) {
		// Counting the number of tools and resources tested

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := client.ListTools(ctx, toolsRequest)
		require.NoError(t, err, "expected to discover toolssuccessfully")
		t.Logf("Discovered %d tools", len(toolsResult.Tools))
		require.Equal(t, 85, len(toolsResult.Tools), "expected to discover 85 tools")
	})

	t.Run("Resource count Tests", func(t *testing.T) {
		// Counting the number of tools and resources tested

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		resourceRequest := mcp.ListResourcesRequest{}
		resourceResult, err := client.ListResources(ctx, resourceRequest)
		require.NoError(t, err, "expected to discover toolssuccessfully")
		t.Logf("Discovered %d resources", len(resourceResult.Resources))
		require.Equal(t, 1, len(resourceResult.Resources), "expected to discover 1 resource")
	})

}

// createStdioClient creates a stdio-based MCP client
func createStdioClient(t *testing.T) (mcpClient.MCPClient, func()) {
	args := []string{
		"docker",
		"run",
		"-i",
		"--rm",
		"-e", "CONSUL_MCP_SERVER_READ_GITHUB_RESOURCES=false",
		"consul-mcp-server:test-e2e",
	}
	t.Log("Starting Stdio MCP client...")
	client, err := mcpClient.NewStdioMCPClient(args[0], []string{}, args[1:]...)
	require.NoError(t, err, "expected to create stdio client successfully")

	cleanup := func() {
		client.Close()
	}

	return client, cleanup
}

// createHTTPClient creates an HTTP-based MCP client
func createHTTPClient(t *testing.T) (mcpClient.MCPClient, func()) {
	t.Log("Starting HTTP MCP server...")

	port := getTestPort()
	baseURL := fmt.Sprintf("http://localhost:%s", port)
	mcpURL := fmt.Sprintf("http://localhost:%s/mcp", port)

	// Start container in HTTP mode
	containerID := startHTTPContainer(t, port)

	// Ensure container cleanup even if test fails
	t.Cleanup(func() {
		stopContainer(t, containerID)
	})

	// Wait for server to be ready
	waitForServer(t, baseURL)

	// Create client with MCP endpoint
	client, err := mcpClient.NewStreamableHttpClient(mcpURL)
	require.NoError(t, err, "expected to create HTTP client successfully")

	cleanup := func() {
		if client != nil {
			client.Close()
		}
		// Container cleanup handled by t.Cleanup()
	}

	return client, cleanup
}

// startHTTPContainer starts a Docker container in HTTP mode and returns container ID
func startHTTPContainer(t *testing.T, port string) string {
	portMapping := fmt.Sprintf("%s:8080", port)
	cmd := exec.Command(
		"docker", "run", "-d", "--rm",
		"-e", "TRANSPORT_MODE=streamable-http",
		"-e", "TRANSPORT_HOST=0.0.0.0",
		"-e", "MCP_SESSION_MODE=stateful",
		"-e", "CONSUL_MCP_SERVER_READ_GITHUB_RESOURCES=false",
		"-p", portMapping,
		"consul-mcp-server:test-e2e",
	)
	output, err := cmd.Output()
	require.NoError(t, err, "expected to start HTTP container successfully")

	containerID := string(output)[:12] // First 12 chars of container ID
	t.Logf("Started HTTP container: %s on port %s", containerID, port)
	return containerID
}

// waitForServer waits for the HTTP server to be ready
func waitForServer(t *testing.T, baseURL string) {
	client := &http.Client{Timeout: 2 * time.Second}
	for i := 0; i < 30; i++ {
		resp, err := client.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			t.Log("HTTP server is ready")
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatal("HTTP server failed to start within 30 seconds")
}

// stopContainer stops the Docker container
func stopContainer(t *testing.T, containerID string) {
	if containerID == "" {
		return
	}

	t.Logf("Stopping container: %s", containerID)
	cmd := exec.Command("docker", "stop", containerID)
	if err := cmd.Run(); err != nil {
		t.Logf("Warning: failed to stop container %s: %v", containerID, err)
		// Try force kill if stop fails
		killCmd := exec.Command("docker", "kill", containerID)
		if killErr := killCmd.Run(); killErr != nil {
			t.Logf("Warning: failed to kill container %s: %v", containerID, killErr)
		}
	} else {
		t.Logf("Successfully stopped container: %s", containerID)
	}
}
