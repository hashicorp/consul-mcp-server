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

func GetCatalogServicesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("catalog_services",
		mcp.WithDescription("Returns the list of services in the Consul catalog."),
		mcp.WithTitleAnnotation("List services in the Consul catalog"),
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
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getCatalogServicesHandler(ctx, request, logger)
		},
	}
}

func getCatalogServicesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
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

	serviceResp, err := consulClient.Get("catalog/services", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching service list from consul catalog", err)
	}

	// convert serviceResp i.e. bytes[] to text
	servicesJson := strings.TrimSpace(string(serviceResp))
	return mcp.NewToolResultText(servicesJson), nil
}

func GetCatalogNodesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("catalog_nodes",
		mcp.WithDescription("Returns the list of nodes in the Consul catalog."),
		mcp.WithTitleAnnotation("List nodes in the Consul catalog"),
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
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getCatalogNodesHandler(ctx, request, logger)
		},
	}
}

func getCatalogNodesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
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

	nodeResp, err := consulClient.Get("catalog/nodes", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching node list from consul catalog", err)
	}

	// convert nodeResp i.e. bytes[] to text
	nodesJson := strings.TrimSpace(string(nodeResp))
	return mcp.NewToolResultText(nodesJson), nil
}

func GetCatalogServiceTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("catalog_service",
		mcp.WithDescription("Returns the list of nodes providing a specific service in the Consul catalog."),
		mcp.WithTitleAnnotation("List nodes for a specific service in the Consul catalog"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_name",
			mcp.Description("The name of the service to query nodes for."),
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
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
		mcp.WithString("tag",
			mcp.Description("Filter results to only nodes providing the service with this tag."),
		),
		mcp.WithString("near",
			mcp.Description("Sort the node list in ascending order based on the estimated round trip time from the given node."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getCatalogServiceHandler(ctx, request, logger)
		},
	}
}

func getCatalogServiceHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceName, err := request.RequireString("service_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")
	tag := request.GetString("tag", "")
	near := request.GetString("near", "")

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

	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if tag != "" {
		queryParams.Set("tag", tag)
	}
	if near != "" {
		queryParams.Set("near", near)
	}

	serviceResp, err := consulClient.Get(fmt.Sprintf("catalog/service/%s", serviceName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching nodes for service '%s' from consul catalog", serviceName), err)
	}

	// convert serviceResp i.e. bytes[] to text
	serviceJson := strings.TrimSpace(string(serviceResp))
	return mcp.NewToolResultText(serviceJson), nil
}

func GetCatalogConnectTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("catalog_connect",
		mcp.WithDescription("Returns the list of nodes for a mesh-capable (Connect-enabled) service in the Consul catalog."),
		mcp.WithTitleAnnotation("List nodes for mesh-capable service in the Consul catalog"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_name",
			mcp.Description("The name of the mesh-capable service to query."),
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
			return getCatalogConnectHandler(ctx, request, logger)
		},
	}
}

func getCatalogConnectHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceName, err := request.RequireString("service_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_name is required", err)
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

	connectResp, err := consulClient.Get(fmt.Sprintf("catalog/connect/%s", serviceName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching mesh-capable service '%s' nodes from consul catalog", serviceName), err)
	}

	// convert connectResp i.e. bytes[] to text
	connectJson := strings.TrimSpace(string(connectResp))
	return mcp.NewToolResultText(connectJson), nil
}

func GetCatalogNodeTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("catalog_node",
		mcp.WithDescription("Returns the services for a specific node in the Consul catalog."),
		mcp.WithTitleAnnotation("List services for a specific node in the Consul catalog"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("node_name",
			mcp.Description("The name of the node to query services for."),
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
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getCatalogNodeHandler(ctx, request, logger)
		},
	}
}

func getCatalogNodeHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	nodeName, err := request.RequireString("node_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: node_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")

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

	if dc != "" {
		queryParams.Set("dc", dc)
	}

	nodeResp, err := consulClient.Get(fmt.Sprintf("catalog/node/%s", nodeName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching services for node '%s' from consul catalog", nodeName), err)
	}

	// convert nodeResp i.e. bytes[] to text
	nodeJson := strings.TrimSpace(string(nodeResp))
	return mcp.NewToolResultText(nodeJson), nil
}

func GetCatalogDatacentersTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("catalog_datacenters",
		mcp.WithDescription("Returns the list of datacenters known to the Consul cluster."),
		mcp.WithTitleAnnotation("List datacenters known to the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getCatalogDatacentersHandler(ctx, request, logger)
		},
	}
}

func getCatalogDatacentersHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	dcResp, err := consulClient.Get("catalog/datacenters", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching datacenters list from consul catalog", err)
	}

	// convert dcResp i.e. bytes[] to text
	dcJson := strings.TrimSpace(string(dcResp))
	return mcp.NewToolResultText(dcJson), nil
}

func GetCatalogGatewayServicesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("catalog_gateway_services",
		mcp.WithDescription("Returns the services associated with a gateway in the Consul catalog."),
		mcp.WithTitleAnnotation("List services for a gateway in the Consul catalog"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("gateway_name",
			mcp.Description("The name of the gateway to query services for."),
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
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getCatalogGatewayServicesHandler(ctx, request, logger)
		},
	}
}

func getCatalogGatewayServicesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	gatewayName, err := request.RequireString("gateway_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: gateway_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")

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

	if dc != "" {
		queryParams.Set("dc", dc)
	}

	gatewayResp, err := consulClient.Get(fmt.Sprintf("catalog/gateway-services/%s", gatewayName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching services for gateway '%s' from consul catalog", gatewayName), err)
	}

	// convert gatewayResp i.e. bytes[] to text
	gatewayJson := strings.TrimSpace(string(gatewayResp))
	return mcp.NewToolResultText(gatewayJson), nil
}
