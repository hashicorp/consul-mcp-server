// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"

	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
)

// contextKey is a type alias to avoid lint warnings while maintaining compatibility
type contextKey string

// NewSessionHandler initializes clients for the session
func NewSessionHandler(ctx context.Context, session server.ClientSession, logger *log.Logger) {
	// Create a unique Consul client per session
	consulClient := NewConsulClient(ctx, session.SessionID(), logger)

	if consulClient != nil {
		logger.WithField("session_id", session.SessionID()).Info("Created Consul HTTP client for session")
	} else {
		logger.WithField("session_id", session.SessionID()).Warn("Consul HTTP client is nil for session")
	}
}

// EndSessionHandler cleans up clients when the session ends
func EndSessionHandler(_ context.Context, session server.ClientSession, logger *log.Logger) {
	DeleteConsulHttpClientForSession(session.SessionID())
	logger.WithField("session_id", session.SessionID()).Info("Cleaned up clients for session")
}
