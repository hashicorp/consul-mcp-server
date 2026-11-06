// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

const initializeRequest = `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}`

// TODO: e2e tests only run the server in `stdio` mode (stdin/stdout), never `streamable-http` mode.
// Add HTTP e2e tests against the real Docker image; these httptest-based tests do not cover TLS or the binary's flags/env.

func newTestHTTPServer(t *testing.T) *httptest.Server {
	t.Helper()
	logger, _ := test.NewNullLogger()
	srv := httptest.NewServer(newStreamableHTTPHandler(NewServer("test", logger), logger, "/mcp", false))
	t.Cleanup(srv.Close)
	return srv
}

func postMCP(t *testing.T, url string, body string, chunked bool) *http.Response {
	t.Helper()
	return postMCPWithHost(t, url, "", body, chunked)
}

func postMCPWithHost(t *testing.T, url, host, body string, chunked bool) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url+"/mcp", strings.NewReader(body))
	require.NoError(t, err)
	if host != "" {
		req.Host = host
	}
	if chunked {
		req.ContentLength = -1
		req.TransferEncoding = []string{"chunked"}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func TestStreamableHTTPRequestBodyLimit(t *testing.T) {
	srv := newTestHTTPServer(t)
	oversized := `{"jsonrpc":"2.0","id":1,"method":"ping","params":{"pad":"` +
		strings.Repeat("a", server.DefaultMaxRequestBodyBytes) + `"}}`

	t.Run("within limit succeeds", func(t *testing.T) {
		resp := postMCP(t, srv.URL, initializeRequest, false)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("oversized with content length is rejected", func(t *testing.T) {
		resp := postMCP(t, srv.URL, oversized, false)
		require.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
	})

	t.Run("oversized chunked is rejected", func(t *testing.T) {
		resp := postMCP(t, srv.URL, oversized, true)
		require.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
	})
}

// httptest listens on loopback, so mcp-go's DNS-rebinding check applies.
func TestStreamableHTTPDNSRebindingProtection(t *testing.T) {
	srv := newTestHTTPServer(t)

	t.Run("localhost host is allowed", func(t *testing.T) {
		resp := postMCPWithHost(t, srv.URL, "localhost", initializeRequest, false)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-localhost host is rejected", func(t *testing.T) {
		resp := postMCPWithHost(t, srv.URL, "evil.example", initializeRequest, false)
		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func TestSlogLoggerPreservesLevel(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(log.DebugLevel)
	slogger := newSlogLogger(logger).WithGroup("transport").With("component", "http")

	slogger.Debug("debug msg")
	slogger.Info("info msg")
	slogger.Warn("warn msg")
	slogger.Error("error msg", "err", "boom")

	entries := hook.AllEntries()
	require.Len(t, entries, 4)
	require.Equal(t, log.DebugLevel, entries[0].Level)
	require.Equal(t, log.InfoLevel, entries[1].Level)
	require.Equal(t, log.WarnLevel, entries[2].Level)
	require.Equal(t, log.ErrorLevel, entries[3].Level)
	require.Equal(t, "error msg", entries[3].Message)
	require.Equal(t, "http", entries[3].Data["transport.component"])
	require.Equal(t, "boom", entries[3].Data["transport.err"])
}

func TestSlogLoggerRespectsLogrusLevel(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(log.InfoLevel)
	slogger := newSlogLogger(logger)

	slogger.Debug("debug msg")
	slogger.Info("info msg")

	entries := hook.AllEntries()
	require.Len(t, entries, 1)
	require.Equal(t, "info msg", entries[0].Message)
}
