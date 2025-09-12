// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/consul-mcp-server/pkg/client"
	"github.com/hashicorp/consul-mcp-server/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetOperatorRaftConfigurationTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("operator_raft_configuration",
		mcp.WithDescription("Returns the current Raft configuration for the Consul cluster."),
		mcp.WithTitleAnnotation("Get Raft configuration from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
		mcp.WithString("stale",
			mcp.Description("If present, results from a stale follower may be used. Set to 'true' to enable."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getOperatorRaftConfigurationHandler(ctx, request, logger)
		},
	}
}

func getOperatorRaftConfigurationHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get optional parameters
	dc := request.GetString("dc", "")
	stale := request.GetString("stale", "")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{}
	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if stale == "true" {
		queryParams.Set("stale", "true")
	}

	uri := (&url.URL{
		Path:     "operator/raft/configuration",
		RawQuery: queryParams.Encode(),
	}).String()

	raftResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching Raft configuration from consul operator", err)
	}

	// convert raftResp i.e. bytes[] to text
	raftJson := strings.TrimSpace(string(raftResp))
	return mcp.NewToolResultText(raftJson), nil
}
