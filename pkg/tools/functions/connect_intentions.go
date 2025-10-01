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

func GetConnectIntentionsTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("connect_intentions",
		mcp.WithDescription("Returns the list of Connect service intentions in the Consul cluster."),
		mcp.WithTitleAnnotation("List Connect service intentions in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("namespace",
			mcp.Description("The consul namespace to query."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("filter",
			mcp.Description("Specifies the expression used to filter the queries results prior to returning the data."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getConnectIntentionsHandler(ctx, request, logger)
		},
	}
}

func getConnectIntentionsHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	filter := request.GetString("filter", "")

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

	if filter != "" {
		queryParams.Set("filter", filter)
	}

	intentionsResp, err := consulClient.Get("connect/intentions", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching Connect intentions from consul", err)
	}

	// convert intentionsResp i.e. bytes[] to text
	intentionsJson := strings.TrimSpace(string(intentionsResp))
	return mcp.NewToolResultText(intentionsJson), nil
}

func GetConnectIntentionTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("connect_intention",
		mcp.WithDescription("Returns the details of a specific Connect service intention by ID."),
		mcp.WithTitleAnnotation("Get specific Connect service intention details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("intention_id",
			mcp.Description("The UUID of the intention to query."),
			mcp.Required(),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getConnectIntentionHandler(ctx, request, logger)
		},
	}
}

func getConnectIntentionHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	intentionId, err := request.RequireString("intention_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: intention_id is required", err)
	}

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	intentionResp, err := consulClient.Get(fmt.Sprintf("connect/intentions/%s", intentionId), nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching Connect intention '%s' from consul", intentionId), err)
	}

	// convert intentionResp i.e. bytes[] to text
	intentionJson := strings.TrimSpace(string(intentionResp))
	return mcp.NewToolResultText(intentionJson), nil
}

func GetConnectIntentionMatchTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("connect_intention_match",
		mcp.WithDescription("Returns the list of intentions that match a given source or destination service."),
		mcp.WithTitleAnnotation("Match Connect intentions by source or destination service"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("by",
			mcp.Description("Specifies whether to match by 'source' or 'destination' service."),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("The name of the service to match intentions for."),
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
			return getConnectIntentionMatchHandler(ctx, request, logger)
		},
	}
}

func getConnectIntentionMatchHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	by, err := request.RequireString("by")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: by is required", err)
	}

	name, err := request.RequireString("name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: name is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Validate 'by' parameter
	if by != "source" && by != "destination" {
		return nil, utils.LogAndReturnError(logger, "invalid 'by' parameter: must be 'source' or 'destination'", fmt.Errorf("invalid by parameter: %s", by))
	}

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{
		"by":        {by},
		"name":      {name},
		"partition": {ap},
		"ns":        {ns},
	}

	matchResp, err := consulClient.Get("connect/intentions/match", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching Connect intention matches for %s '%s' from consul", by, name), err)
	}

	// convert matchResp i.e. bytes[] to text
	matchJson := strings.TrimSpace(string(matchResp))
	return mcp.NewToolResultText(matchJson), nil
}

func GetConnectIntentionCheckTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("connect_intention_check",
		mcp.WithDescription("Checks whether a connection between two services is authorized by Connect intentions."),
		mcp.WithTitleAnnotation("Check if connection is authorized by Connect intentions"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("source",
			mcp.Description("The name of the source service."),
			mcp.Required(),
		),
		mcp.WithString("destination",
			mcp.Description("The name of the destination service."),
			mcp.Required(),
		),
		mcp.WithString("source_type",
			mcp.Description("The type of the source (default: 'consul'). Can be 'consul' for Consul services."),
			mcp.DefaultString("consul"),
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
			return getConnectIntentionCheckHandler(ctx, request, logger)
		},
	}
}

func getConnectIntentionCheckHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	source, err := request.RequireString("source")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: source is required", err)
	}

	destination, err := request.RequireString("destination")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: destination is required", err)
	}

	sourceType := request.GetString("source_type", "consul")
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
		"source":      {source},
		"destination": {destination},
		"source-type": {sourceType},
		"partition":   {ap},
		"ns":          {ns},
	}

	checkResp, err := consulClient.Get("connect/intentions/check", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("checking Connect intention authorization from '%s' to '%s' from consul", source, destination), err)
	}

	// convert checkResp i.e. bytes[] to text
	checkJson := strings.TrimSpace(string(checkResp))
	return mcp.NewToolResultText(checkJson), nil
}
