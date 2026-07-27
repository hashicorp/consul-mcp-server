// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConsulContextMiddlewareStripsSessionIDInStatelessMode(t *testing.T) {
	logger := log.New()
	handler := ConsulContextMiddleware(logger, true)(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		assert.True(t, IsStatelessRequest(r.Context()))
		assert.Empty(t, r.Header.Get(server.HeaderKeySessionID))

		sessionID, ok := contextStringValue(r.Context(), contextKeySessionID)
		assert.False(t, ok)
		assert.Empty(t, sessionID)
	}))

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set(server.HeaderKeySessionID, "attacker-session")

	handler.ServeHTTP(httptest.NewRecorder(), req)
}

func TestConsulContextMiddlewareKeepsSessionIDInStatefulMode(t *testing.T) {
	logger := log.New()
	handler := ConsulContextMiddleware(logger, false)(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		assert.False(t, IsStatelessRequest(r.Context()))
		assert.Equal(t, "server-session", r.Header.Get(server.HeaderKeySessionID))

		sessionID, ok := contextStringValue(r.Context(), contextKeySessionID)
		assert.True(t, ok)
		assert.Equal(t, "server-session", sessionID)
	}))

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set(server.HeaderKeySessionID, "server-session")

	handler.ServeHTTP(httptest.NewRecorder(), req)
}
