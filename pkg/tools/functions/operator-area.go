// Copyright (c) HashiCorp, Inc.
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

func GetOperatorAreasTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("operator_areas",
		mcp.WithDescription("Returns the list of network areas in the Consul cluster."),
		mcp.WithTitleAnnotation("List network areas in the Consul cluster"),
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
			return getOperatorAreasHandler(ctx, request, logger)
		},
	}
}

func getOperatorAreasHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
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

	areasResp, err := consulClient.Get("operator/area", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching network areas from consul operator", err)
	}

	// convert areasResp i.e. bytes[] to text
	areasJson := strings.TrimSpace(string(areasResp))
	return mcp.NewToolResultText(areasJson), nil
}

func GetOperatorAreaTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("operator_area",
		mcp.WithDescription("Returns the details of a specific network area in the Consul cluster."),
		mcp.WithTitleAnnotation("Get specific network area details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("area_id",
			mcp.Description("The UUID of the network area to query."),
			mcp.Required(),
		),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getOperatorAreaHandler(ctx, request, logger)
		},
	}
}

func getOperatorAreaHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	areaId, err := request.RequireString("area_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: area_id is required", err)
	}

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

	areaResp, err := consulClient.Get(fmt.Sprintf("operator/area/%s", areaId), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching network area '%s' from consul operator", areaId), err)
	}

	// convert areaResp i.e. bytes[] to text
	areaJson := strings.TrimSpace(string(areaResp))
	return mcp.NewToolResultText(areaJson), nil
}

func GetOperatorAreaMembersTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("operator_area_members",
		mcp.WithDescription("Returns the list of members in a specific network area."),
		mcp.WithTitleAnnotation("List members in a specific network area"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("area_id",
			mcp.Description("The UUID of the network area to query members for."),
			mcp.Required(),
		),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getOperatorAreaMembersHandler(ctx, request, logger)
		},
	}
}

func getOperatorAreaMembersHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	areaId, err := request.RequireString("area_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: area_id is required", err)
	}

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

	membersResp, err := consulClient.Get(fmt.Sprintf("operator/area/%s/members", areaId), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching members for network area '%s' from consul operator", areaId), err)
	}

	// convert membersResp i.e. bytes[] to text
	membersJson := strings.TrimSpace(string(membersResp))
	return mcp.NewToolResultText(membersJson), nil
}
