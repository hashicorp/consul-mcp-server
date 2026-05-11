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

func GetConfigEntriesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("config_entries",
		mcp.WithDescription("Returns the list of configuration entries of a specific kind in the Consul cluster."),
		mcp.WithTitleAnnotation("List configuration entries of a specific kind in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("kind",
			mcp.Description("The kind of configuration entry to query (e.g., service-defaults, proxy-defaults, service-router, etc.)."),
			mcp.Required(),
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
			return getConfigEntriesHandler(ctx, request, logger)
		},
	}
}

func getConfigEntriesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	kind, err := request.RequireString("kind")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: kind is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{
		"partition": {ap},
		"ns":        {ns},
	}

	configResp, err := consulClient.Get(fmt.Sprintf("config/%s", kind), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching config entries of kind '%s' from consul", kind), err)
	}

	// convert configResp i.e. bytes[] to text
	configJson := strings.TrimSpace(string(configResp))
	return mcp.NewToolResultText(configJson), nil
}

func GetConfigEntryTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("config_entry",
		mcp.WithDescription("Returns the details of a specific configuration entry in the Consul cluster."),
		mcp.WithTitleAnnotation("Get specific configuration entry details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("kind",
			mcp.Description("The kind of configuration entry (e.g., service-defaults, proxy-defaults, service-router, etc.)."),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("The name of the configuration entry to query."),
			mcp.Required(),
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
			return getConfigEntryHandler(ctx, request, logger)
		},
	}
}

func getConfigEntryHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	kind, err := request.RequireString("kind")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: kind is required", err)
	}

	name, err := request.RequireString("name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: name is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{
		"partition": {ap},
		"ns":        {ns},
	}

	configResp, err := consulClient.Get(fmt.Sprintf("config/%s/%s", kind, name), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching config entry '%s' of kind '%s' from consul", name, kind), err)
	}

	// convert configResp i.e. bytes[] to text
	configJson := strings.TrimSpace(string(configResp))
	return mcp.NewToolResultText(configJson), nil
}

func GetConfigKindsTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("config_kinds",
		mcp.WithDescription("Returns the list of configuration entry kinds available in the Consul cluster."),
		mcp.WithTitleAnnotation("List configuration entry kinds in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getConfigKindsHandler(ctx, request, logger)
		},
	}
}

func getConfigKindsHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

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

	kindsResp, err := consulClient.Get("config", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching config entry kinds from consul", err)
	}

	// convert kindsResp i.e. bytes[] to text
	kindsJson := strings.TrimSpace(string(kindsResp))
	return mcp.NewToolResultText(kindsJson), nil
}
