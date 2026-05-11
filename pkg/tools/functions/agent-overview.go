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

func GetAgentSelfTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_self",
		mcp.WithDescription("Returns the local agent's configuration and member information."),
		mcp.WithTitleAnnotation("Get local agent configuration and member information"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentSelfHandler(ctx, request, logger)
		},
	}
}

func getAgentSelfHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	selfResp, err := consulClient.Get("agent/self", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching agent self information", err)
	}

	selfJson := strings.TrimSpace(string(selfResp))
	return mcp.NewToolResultText(selfJson), nil
}

func GetAgentConfigTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_config",
		mcp.WithDescription("Returns the configuration of the local agent."),
		mcp.WithTitleAnnotation("Get local agent configuration"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentConfigHandler(ctx, request, logger)
		},
	}
}

func getAgentConfigHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	configResp, err := consulClient.Get("agent/config", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching agent configuration", err)
	}

	configJson := strings.TrimSpace(string(configResp))
	return mcp.NewToolResultText(configJson), nil
}

func GetAgentMembersTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_members",
		mcp.WithDescription("Returns the members the agent sees in the cluster gossip pool."),
		mcp.WithTitleAnnotation("Get cluster members seen by the agent"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("wan",
			mcp.Description("Set to '1' to get WAN pool members instead of LAN pool members."),
		),
		mcp.WithString("segment",
			mcp.Description("Filter results to nodes in the given segment."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentMembersHandler(ctx, request, logger)
		},
	}
}

func getAgentMembersHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	wan := request.GetString("wan", "")
	segment := request.GetString("segment", "")

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{}
	if wan == "true" {
		queryParams.Set("wan", "1")
	}
	if segment != "" {
		queryParams.Set("segment", segment)
	}

	membersResp, err := consulClient.Get("agent/members", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching agent members", err)
	}

	membersJson := strings.TrimSpace(string(membersResp))
	return mcp.NewToolResultText(membersJson), nil
}

func GetAgentMetricsTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_metrics",
		mcp.WithDescription("Returns the current metrics for the agent."),
		mcp.WithTitleAnnotation("Get current agent metrics"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("format",
			mcp.Description("Format for the metrics output (prometheus or JSON)."),
			mcp.DefaultString("json"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentMetricsHandler(ctx, request, logger)
		},
	}
}

func getAgentMetricsHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	format := request.GetString("format", "")

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{}
	if format != "" {
		queryParams.Set("format", format)
	}

	metricsResp, err := consulClient.Get("agent/metrics", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching agent metrics", err)
	}

	metricsOutput := strings.TrimSpace(string(metricsResp))
	return mcp.NewToolResultText(metricsOutput), nil
}

func GetAgentHostTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_host",
		mcp.WithDescription("Returns information about the host the agent is running on."),
		mcp.WithTitleAnnotation("Get agent host information"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentHostHandler(ctx, request, logger)
		},
	}
}

func getAgentHostHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	hostResp, err := consulClient.Get("agent/host", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching agent host information", err)
	}

	hostJson := strings.TrimSpace(string(hostResp))
	return mcp.NewToolResultText(hostJson), nil
}

func GetAgentVersionTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_version",
		mcp.WithDescription("Returns the version of the local agent."),
		mcp.WithTitleAnnotation("Get local agent version"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentVersionHandler(ctx, request, logger)
		},
	}
}

func getAgentVersionHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	versionResp, err := consulClient.Get("agent/version", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching agent version", err)
	}

	versionJson := strings.TrimSpace(string(versionResp))
	return mcp.NewToolResultText(versionJson), nil
}

func GetAgentReloadTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_reload",
		mcp.WithDescription("Triggers the agent to reload its configuration."),
		mcp.WithTitleAnnotation("Reload agent configuration"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentReloadHandler(ctx, request, logger)
		},
	}
}

func getAgentReloadHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	reloadResp, err := consulClient.Get("agent/reload", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "triggering agent reload", err)
	}

	reloadJson := strings.TrimSpace(string(reloadResp))
	return mcp.NewToolResultText(reloadJson), nil
}
