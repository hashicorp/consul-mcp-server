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

func GetAgentServicesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_services",
		mcp.WithDescription("Returns all services that are registered with the local agent."),
		mcp.WithTitleAnnotation("List all services registered with the local agent"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("filter",
			mcp.Description("Filter expression to use for filtering the services."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentServicesHandler(ctx, request, logger)
		},
	}
}

func getAgentServicesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
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

	servicesResp, err := consulClient.Get("agent/services", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching agent services", err)
	}

	servicesJson := strings.TrimSpace(string(servicesResp))
	return mcp.NewToolResultText(servicesJson), nil
}

func GetAgentServiceTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_service",
		mcp.WithDescription("Returns the service definition for a specific service registered with the local agent."),
		mcp.WithTitleAnnotation("Get a specific service definition by ID"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_id",
			mcp.Description("The ID of the service to retrieve."),
			mcp.Required(),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentServiceHandler(ctx, request, logger)
		},
	}
}

func getAgentServiceHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceID, err := request.RequireString("service_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_id is required", err)
	}

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	serviceResp, err := consulClient.Get(fmt.Sprintf("agent/service/%s", serviceID), nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching agent service '%s'", serviceID), err)
	}

	serviceJson := strings.TrimSpace(string(serviceResp))
	return mcp.NewToolResultText(serviceJson), nil
}

func GetAgentServiceConfigurationTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_service_configuration",
		mcp.WithDescription("Returns the full service definition for a specific service registered with the local agent."),
		mcp.WithTitleAnnotation("Get full service configuration by ID"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_id",
			mcp.Description("The ID of the service to retrieve configuration for."),
			mcp.Required(),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentServiceConfigurationHandler(ctx, request, logger)
		},
	}
}

func getAgentServiceConfigurationHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceID, err := request.RequireString("service_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_id is required", err)
	}

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	configResp, err := consulClient.Get(fmt.Sprintf("agent/service/%s/configuration", serviceID), nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching agent service configuration for '%s'", serviceID), err)
	}

	configJson := strings.TrimSpace(string(configResp))
	return mcp.NewToolResultText(configJson), nil
}

func GetAgentServiceLocalStateTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_service_local_state",
		mcp.WithDescription("Returns the local state for all services registered with the agent."),
		mcp.WithTitleAnnotation("Get local state for all services"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentServiceLocalStateHandler(ctx, request, logger)
		},
	}
}

func getAgentServiceLocalStateHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	stateResp, err := consulClient.Get("agent/local-state", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching agent local state", err)
	}

	stateJson := strings.TrimSpace(string(stateResp))
	return mcp.NewToolResultText(stateJson), nil
}
