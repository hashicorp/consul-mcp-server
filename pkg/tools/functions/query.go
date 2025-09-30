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

func GetQueryTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("query",
		mcp.WithDescription("Returns the list of prepared queries in the Consul cluster."),
		mcp.WithTitleAnnotation("List prepared queries in the Consul cluster"),
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
			return getQueryHandler(ctx, request, logger)
		},
	}
}

func getQueryHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
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

	queryResp, err := consulClient.Get("query", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching prepared queries from consul", err)
	}

	// convert queryResp i.e. bytes[] to text
	queryJson := strings.TrimSpace(string(queryResp))
	return mcp.NewToolResultText(queryJson), nil
}

func GetQueryByIdTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("query_by_id",
		mcp.WithDescription("Returns the details of a specific prepared query by ID."),
		mcp.WithTitleAnnotation("Get specific prepared query details by ID"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("query_id",
			mcp.Description("The UUID of the prepared query to retrieve."),
			mcp.Required(),
		),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getQueryByIdHandler(ctx, request, logger)
		},
	}
}

func getQueryByIdHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	queryId, err := request.RequireString("query_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: query_id is required", err)
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

	queryResp, err := consulClient.Get(fmt.Sprintf("query/%s", queryId), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching prepared query '%s' from consul", queryId), err)
	}

	// convert queryResp i.e. bytes[] to text
	queryJson := strings.TrimSpace(string(queryResp))
	return mcp.NewToolResultText(queryJson), nil
}

func GetQueryExecuteTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("query_execute",
		mcp.WithDescription("Executes a prepared query and returns the results."),
		mcp.WithTitleAnnotation("Execute a prepared query and get results"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("query_id_or_name",
			mcp.Description("The UUID or name of the prepared query to execute."),
			mcp.Required(),
		),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
		mcp.WithString("near",
			mcp.Description("Sort the results by network coordinate distance from the given node."),
		),
		mcp.WithString("limit",
			mcp.Description("Limit the number of results returned."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getQueryExecuteHandler(ctx, request, logger)
		},
	}
}

func getQueryExecuteHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	queryIdOrName, err := request.RequireString("query_id_or_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: query_id_or_name is required", err)
	}

	// Get optional parameters
	dc := request.GetString("dc", "")
	near := request.GetString("near", "")
	limit := request.GetString("limit", "")

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
	if near != "" {
		queryParams.Set("near", near)
	}
	if limit != "" {
		queryParams.Set("limit", limit)
	}

	executeResp, err := consulClient.Get(fmt.Sprintf("query/%s/execute", queryIdOrName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("executing prepared query '%s' from consul", queryIdOrName), err)
	}

	// convert executeResp i.e. bytes[] to text
	executeJson := strings.TrimSpace(string(executeResp))
	return mcp.NewToolResultText(executeJson), nil
}

func GetQueryExplainTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("query_explain",
		mcp.WithDescription("Explains how a prepared query will be executed without running it."),
		mcp.WithTitleAnnotation("Explain how a prepared query will be executed"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("query_id_or_name",
			mcp.Description("The UUID or name of the prepared query to explain."),
			mcp.Required(),
		),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getQueryExplainHandler(ctx, request, logger)
		},
	}
}

func getQueryExplainHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	queryIdOrName, err := request.RequireString("query_id_or_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: query_id_or_name is required", err)
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

	explainResp, err := consulClient.Get(fmt.Sprintf("query/%s/explain", queryIdOrName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("explaining prepared query '%s' from consul", queryIdOrName), err)
	}

	// convert explainResp i.e. bytes[] to text
	explainJson := strings.TrimSpace(string(explainResp))
	return mcp.NewToolResultText(explainJson), nil
}
