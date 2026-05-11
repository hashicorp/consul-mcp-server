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

func GetAgentChecksTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_checks",
		mcp.WithDescription("Returns all checks that are registered with the local agent."),
		mcp.WithTitleAnnotation("List all checks registered with the local agent"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("filter",
			mcp.Description("Filter expression to use for filtering the checks."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentChecksHandler(ctx, request, logger)
		},
	}
}

func getAgentChecksHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	filter := request.GetString("filter", "")

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{}
	if filter != "" {
		queryParams.Set("filter", filter)
	}

	checksResp, err := consulClient.Get("agent/checks", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching agent checks", err)
	}

	checksJson := strings.TrimSpace(string(checksResp))
	return mcp.NewToolResultText(checksJson), nil
}

func GetAgentCheckTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_check",
		mcp.WithDescription("Returns the health check specified by the given ID."),
		mcp.WithTitleAnnotation("Get a specific health check by ID"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("check_id",
			mcp.Description("The ID of the health check to retrieve."),
			mcp.Required(),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentCheckHandler(ctx, request, logger)
		},
	}
}

func getAgentCheckHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	checkID, err := request.RequireString("check_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: check_id is required", err)
	}

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	checkResp, err := consulClient.Get(fmt.Sprintf("agent/health/check/id/%s", checkID), nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching agent check '%s'", checkID), err)
	}

	checkJson := strings.TrimSpace(string(checkResp))
	return mcp.NewToolResultText(checkJson), nil
}

func GetAgentCheckByNameTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_check_by_name",
		mcp.WithDescription("Returns the health check specified by the given name."),
		mcp.WithTitleAnnotation("Get a specific health check by name"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("check_name",
			mcp.Description("The name of the health check to retrieve."),
			mcp.Required(),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentCheckByNameHandler(ctx, request, logger)
		},
	}
}

func getAgentCheckByNameHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	checkName, err := request.RequireString("check_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: check_name is required", err)
	}

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	checkResp, err := consulClient.Get(fmt.Sprintf("agent/health/check/name/%s", checkName), nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching agent check by name '%s'", checkName), err)
	}

	checkJson := strings.TrimSpace(string(checkResp))
	return mcp.NewToolResultText(checkJson), nil
}

func GetAgentServiceHealthTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_service_health",
		mcp.WithDescription("Returns the aggregated health status of a service by ID."),
		mcp.WithTitleAnnotation("Get aggregated health status of a service"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_id",
			mcp.Description("The ID of the service to get health status for."),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Format of the response (text or json)."),
			mcp.DefaultString("json"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentServiceHealthHandler(ctx, request, logger)
		},
	}
}

func getAgentServiceHealthHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceID, err := request.RequireString("service_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_id is required", err)
	}

	format, err := request.RequireString("format")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: format is required", err)
	}

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{
		"format": {format},
	}

	healthResp, err := consulClient.Get(fmt.Sprintf("agent/health/service/id/%s", serviceID), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching health status for service '%s'", serviceID), err)
	}

	healthOutput := strings.TrimSpace(string(healthResp))
	return mcp.NewToolResultText(healthOutput), nil
}

func GetAgentServiceHealthByNameTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_service_health_by_name",
		mcp.WithDescription("Returns the aggregated health status of a service by name."),
		mcp.WithTitleAnnotation("Get aggregated health status of a service by name"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_name",
			mcp.Description("The name of the service to get health status for."),
			mcp.Required(),
		),
		mcp.WithString("format",
			mcp.Description("Format of the response (text or json)."),
			mcp.DefaultString("json"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentServiceHealthByNameHandler(ctx, request, logger)
		},
	}
}

func getAgentServiceHealthByNameHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceName, err := request.RequireString("service_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_name is required", err)
	}

	format, err := request.RequireString("format")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: format is required", err)
	}

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{
		"format": {format},
	}

	healthResp, err := consulClient.Get(fmt.Sprintf("agent/health/service/name/%s", serviceName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching health status for service by name '%s'", serviceName), err)
	}

	healthOutput := strings.TrimSpace(string(healthResp))
	return mcp.NewToolResultText(healthOutput), nil
}
