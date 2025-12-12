// Copyright IBM Corp. 2025
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

func GetSessionTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("session",
		mcp.WithDescription("Returns the details of a specific session in the Consul cluster."),
		mcp.WithTitleAnnotation("Get specific session details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("session_id",
			mcp.Description("The UUID of the session to query."),
			mcp.Required(),
		),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("namespace",
			mcp.Description("The consul namespace to query."),
			mcp.DefaultString("default"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getSessionHandler(ctx, request, logger)
		},
	}
}

func getSessionHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	sessionID, err := request.RequireString("session_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: session_id is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{
		"partition": {ap},
		"ns":        {ns},
	}

	if dc != "" {
		queryParams.Set("dc", dc)
	}

	sessionResp, err := consulClient.Get(fmt.Sprintf("session/info/%s", sessionID), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching session '%s' details from consul", sessionID), err)
	}

	// convert sessionResp i.e. bytes[] to text
	sessionJson := strings.TrimSpace(string(sessionResp))
	return mcp.NewToolResultText(sessionJson), nil
}

func GetSessionNodeTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("session_node",
		mcp.WithDescription("Returns the list of sessions for a specific node."),
		mcp.WithTitleAnnotation("List sessions for a specific node"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("node_name",
			mcp.Description("The name of the node to query sessions for."),
			mcp.Required(),
		),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("namespace",
			mcp.Description("The consul namespace to query."),
			mcp.DefaultString("default"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getSessionNodeHandler(ctx, request, logger)
		},
	}
}

func getSessionNodeHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	node, err := request.RequireString("node")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: node is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{
		"partition": {ap},
		"ns":        {ns},
	}

	if dc != "" {
		queryParams.Set("dc", dc)
	}

	nodeSessionsResp, err := consulClient.Get(fmt.Sprintf("session/node/%s", node), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching sessions for node '%s' from consul", node), err)
	}

	// convert nodeSessionsResp i.e. bytes[] to text
	nodeSessionsJson := strings.TrimSpace(string(nodeSessionsResp))
	return mcp.NewToolResultText(nodeSessionsJson), nil
}

func GetSessionListTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("session_list",
		mcp.WithDescription("Returns the list of all sessions in the Consul cluster."),
		mcp.WithTitleAnnotation("List all sessions in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("namespace",
			mcp.Description("The consul namespace to query."),
			mcp.DefaultString("default"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getSessionListHandler(ctx, request, logger)
		},
	}
}

func getSessionListHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{
		"partition": {ap},
		"ns":        {ns},
	}

	if dc != "" {
		queryParams.Set("dc", dc)
	}

	sessionsResp, err := consulClient.Get("session/list", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching sessions list from consul", err)
	}

	// convert sessionsResp i.e. bytes[] to text
	sessionsJson := strings.TrimSpace(string(sessionsResp))
	return mcp.NewToolResultText(sessionsJson), nil
}
