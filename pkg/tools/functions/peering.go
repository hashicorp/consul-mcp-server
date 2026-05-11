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

func GetPeeringsTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("peerings",
		mcp.WithDescription("Returns the list of peering connections in the Consul cluster."),
		mcp.WithTitleAnnotation("List peering connections in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getPeeringsHandler(ctx, request, logger)
		},
	}
}

func getPeeringsHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	ap := request.GetString("admin_partition", "default")
	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{
		"partition": {ap},
	}
	peeringResp, err := consulClient.Get("peerings", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching peerings list from consul", err)
	}

	// convert peeringResp i.e. bytes[] to text
	peeringsJson := strings.TrimSpace(string(peeringResp))
	return mcp.NewToolResultText(peeringsJson), nil
}

func GetPeeringTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("peering",
		mcp.WithDescription("Returns the details of a specific peering connection in the Consul cluster."),
		mcp.WithTitleAnnotation("Get specific peering connection details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("peering_name",
			mcp.Description("The name of the peering connection to query."),
			mcp.Required(),
		),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getPeeringHandler(ctx, request, logger)
		},
	}
}

func getPeeringHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	peeringName, err := request.RequireString("peering_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: peering_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{
		"partition": {ap},
	}
	peeringResp, err := consulClient.Get(fmt.Sprintf("peering/%s", peeringName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching peering '%s' details from consul", peeringName), err)
	}

	// convert peeringResp i.e. bytes[] to text
	peeringJson := strings.TrimSpace(string(peeringResp))
	return mcp.NewToolResultText(peeringJson), nil
}

func GetPeeringExportedServicesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("peering_exported_services",
		mcp.WithDescription("Returns the list of services exported to a specific peering connection."),
		mcp.WithTitleAnnotation("List services exported to a peering connection"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("peering_name",
			mcp.Description("The name of the peering connection to query exported services for."),
			mcp.Required(),
		),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getPeeringExportedServicesHandler(ctx, request, logger)
		},
	}
}

func getPeeringExportedServicesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	peeringName, err := request.RequireString("peering_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: peering_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{
		"partition": {ap},
	}
	exportedResp, err := consulClient.Get(fmt.Sprintf("peering/%s/exported-services", peeringName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching exported services for peering '%s' from consul", peeringName), err)
	}

	// convert exportedResp i.e. bytes[] to text
	exportedJson := strings.TrimSpace(string(exportedResp))
	return mcp.NewToolResultText(exportedJson), nil
}
