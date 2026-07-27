// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"

	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
)

// contextKey is a type alias to avoid lint warnings while maintaining compatibility
type contextKey string

const (
	contextKeySessionID    contextKey = "session_id"
	contextKeyConsulToken  contextKey = "consul_token"
	contextKeyMCPStateless contextKey = "mcp_stateless"
)

func contextStringValue(ctx context.Context, key contextKey) (string, bool) {
	if value, ok := ctx.Value(key).(string); ok {
		return value, true
	}

	// Preserve compatibility with existing tests and callers that use string keys.
	if value, ok := ctx.Value(string(key)).(string); ok {
		return value, true
	}

	return "", false
}

func contextBoolValue(ctx context.Context, key contextKey) bool {
	if value, ok := ctx.Value(key).(bool); ok {
		return value
	}

	if value, ok := ctx.Value(string(key)).(bool); ok {
		return value
	}

	return false
}

// IsStatelessRequest returns true when this request is being handled by the
// stateless StreamableHTTP transport mode.
func IsStatelessRequest(ctx context.Context) bool {
	return contextBoolValue(ctx, contextKeyMCPStateless)
}

// NewSessionHandler initializes clients for the session
func NewSessionHandler(ctx context.Context, session server.ClientSession, logger *log.Logger) {
	// Create a unique Consul client per session
	consulClient := NewConsulClient(ctx, session.SessionID(), logger)

	if consulClient != nil {
		logger.Info("Created Consul HTTP client for MCP session")
	} else {
		logger.Warn("Consul HTTP client is nil for MCP session")
	}
}

// EndSessionHandler cleans up clients when the session ends
func EndSessionHandler(_ context.Context, session server.ClientSession, logger *log.Logger) {
	DeleteConsulHttpClientForSession(session.SessionID())
	logger.Info("Cleaned up clients for MCP session")
}
