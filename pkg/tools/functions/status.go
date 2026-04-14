// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package functions

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/consul-mcp-server/pkg/client"
	"github.com/hashicorp/consul-mcp-server/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetStatusLeaderTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("status_leader",
		mcp.WithDescription("Returns the current Raft leader for the Consul cluster."),
		mcp.WithTitleAnnotation("Get current Raft leader from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getStatusLeaderHandler(ctx, request, logger)
		},
	}
}

func getStatusLeaderHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	leaderResp, err := consulClient.Get("status/leader", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching Raft leader from consul status", err)
	}

	// convert leaderResp i.e. bytes[] to text
	leaderText := strings.TrimSpace(string(leaderResp))
	return mcp.NewToolResultText(leaderText), nil
}

func GetStatusPeersTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("status_peers",
		mcp.WithDescription("Returns the current list of Raft peers for the Consul cluster."),
		mcp.WithTitleAnnotation("Get current Raft peers from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getStatusPeersHandler(ctx, request, logger)
		},
	}
}

func getStatusPeersHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	peersResp, err := consulClient.Get("status/peers", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching Raft peers from consul status", err)
	}

	// convert peersResp i.e. bytes[] to text
	peersJson := strings.TrimSpace(string(peersResp))
	return mcp.NewToolResultText(peersJson), nil
}
