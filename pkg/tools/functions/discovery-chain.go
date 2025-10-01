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

func GetDiscoveryChainTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("discovery_chain",
		mcp.WithDescription("Returns the compiled discovery chain for a service in the Consul cluster."),
		mcp.WithTitleAnnotation("Get compiled discovery chain for a service in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_name",
			mcp.Description("The name of the service to get the discovery chain for."),
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
		mcp.WithString("datacenter",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
		mcp.WithString("compile_dc",
			mcp.Description("The datacenter to compile the discovery chain for."),
		),
		mcp.WithString("override_mesh_gateway_mode",
			mcp.Description("Override the mesh gateway mode for the discovery chain compilation."),
		),
		mcp.WithString("override_protocol",
			mcp.Description("Override the protocol for the discovery chain compilation."),
		),
		mcp.WithString("override_connect_timeout",
			mcp.Description("Override the connect timeout for the discovery chain compilation."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getDiscoveryChainHandler(ctx, request, logger)
		},
	}
}

func getDiscoveryChainHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceName, err := request.RequireString("service_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	datacenter := request.GetString("datacenter", "")
	compileDc := request.GetString("compile_dc", "")
	overrideMeshGatewayMode := request.GetString("override_mesh_gateway_mode", "")
	overrideProtocol := request.GetString("override_protocol", "")
	overrideConnectTimeout := request.GetString("override_connect_timeout", "")

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

	if datacenter != "" {
		queryParams.Set("datacenter", datacenter)
	}
	if compileDc != "" {
		queryParams.Set("compile-dc", compileDc)
	}
	if overrideMeshGatewayMode != "" {
		queryParams.Set("override-mesh-gateway.mode", overrideMeshGatewayMode)
	}
	if overrideProtocol != "" {
		queryParams.Set("override-protocol", overrideProtocol)
	}
	if overrideConnectTimeout != "" {
		queryParams.Set("override-connect-timeout", overrideConnectTimeout)
	}

	chainResp, err := consulClient.Get(fmt.Sprintf("discovery-chain/%s", serviceName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching discovery chain for service '%s' from consul", serviceName), err)
	}

	// convert chainResp i.e. bytes[] to text
	chainJson := strings.TrimSpace(string(chainResp))
	return mcp.NewToolResultText(chainJson), nil
}
