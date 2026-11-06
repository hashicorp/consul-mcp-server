// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package e2e

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/require"
)

// TestRequestBodyLimitE2E guards the CVE-2026-108859 fix against the real Docker image.
func TestRequestBodyLimitE2E(t *testing.T) {
	buildDockerImage(t)

	port := "8084"
	baseURL := fmt.Sprintf("http://localhost:%s", port)
	mcpURL := baseURL + "/mcp"

	containerID := startHTTPContainerWithCORS(t, port, "disabled", "")
	t.Cleanup(func() {
		_ = exec.Command("docker", "stop", containerID).Run()
	})
	waitForCORSServer(t, baseURL)

	oversized := `{"jsonrpc":"2.0","id":1,"method":"ping","params":{"pad":"` +
		strings.Repeat("a", server.DefaultMaxRequestBodyBytes) + `"}}`

	t.Run("within limit succeeds", func(t *testing.T) {
		sessionID := initializeMCPSession(t, mcpURL, "")
		require.NotEmpty(t, sessionID)
	})

	t.Run("oversized with content length is rejected", func(t *testing.T) {
		require.Equal(t, http.StatusRequestEntityTooLarge, postBody(t, mcpURL, oversized, false))
	})

	t.Run("oversized chunked is rejected", func(t *testing.T) {
		require.Equal(t, http.StatusRequestEntityTooLarge, postBody(t, mcpURL, oversized, true))
	})
}

func postBody(t *testing.T, url, body string, chunked bool) int {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	require.NoError(t, err)
	if chunked {
		req.ContentLength = -1
		req.TransferEncoding = []string{"chunked"}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	return resp.StatusCode
}
