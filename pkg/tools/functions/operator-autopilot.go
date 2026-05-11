// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

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

func GetOperatorAutopilotConfigurationTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("operator_autopilot_configuration",
		mcp.WithDescription("Returns the current Autopilot configuration for the Consul cluster."),
		mcp.WithTitleAnnotation("Get Autopilot configuration from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getOperatorAutopilotConfigurationHandler(ctx, request, logger)
		},
	}
}

func getOperatorAutopilotConfigurationHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get optional parameters
	dc := request.GetString("dc", "")

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

	configResp, err := consulClient.Get("operator/autopilot/configuration", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching Autopilot configuration from consul operator", err)
	}

	// convert configResp i.e. bytes[] to text
	configJson := strings.TrimSpace(string(configResp))
	return mcp.NewToolResultText(configJson), nil
}

func GetOperatorAutopilotHealthTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("operator_autopilot_health",
		mcp.WithDescription("Returns the health status of servers as seen by Autopilot."),
		mcp.WithTitleAnnotation("Get Autopilot health status from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getOperatorAutopilotHealthHandler(ctx, request, logger)
		},
	}
}

func getOperatorAutopilotHealthHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get optional parameters
	dc := request.GetString("dc", "")

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

	healthResp, err := consulClient.Get("operator/autopilot/health", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching Autopilot health from consul operator", err)
	}

	// convert healthResp i.e. bytes[] to text
	healthJson := strings.TrimSpace(string(healthResp))
	return mcp.NewToolResultText(healthJson), nil
}

func GetOperatorAutopilotStateTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("operator_autopilot_state",
		mcp.WithDescription("Returns the state of the Autopilot including servers and voters."),
		mcp.WithTitleAnnotation("Get Autopilot state from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getOperatorAutopilotStateHandler(ctx, request, logger)
		},
	}
}

func getOperatorAutopilotStateHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get optional parameters
	dc := request.GetString("dc", "")

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

	stateResp, err := consulClient.Get("operator/autopilot/state", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching Autopilot state from consul operator", err)
	}

	// convert stateResp i.e. bytes[] to text
	stateJson := strings.TrimSpace(string(stateResp))
	return mcp.NewToolResultText(stateJson), nil
}
