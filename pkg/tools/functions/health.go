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

func GetHealthNodeTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("health_node",
		mcp.WithDescription("Returns the health information for a specific node in the Consul cluster."),
		mcp.WithTitleAnnotation("Get health information for a specific node in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("node_name",
			mcp.Description("The name of the node to query health information for."),
			mcp.Required(),
		),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
		mcp.WithString("filter",
			mcp.Description("Specifies the expression used to filter the queries results prior to returning the data."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getHealthNodeHandler(ctx, request, logger)
		},
	}
}

func getHealthNodeHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	nodeName, err := request.RequireString("node_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: node_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")
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
	}

	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if filter != "" {
		queryParams.Set("filter", filter)
	}

	healthResp, err := consulClient.Get(fmt.Sprintf("health/node/%s", nodeName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching health information for node '%s' from consul", nodeName), err)
	}

	// convert healthResp i.e. bytes[] to text
	healthJson := strings.TrimSpace(string(healthResp))
	return mcp.NewToolResultText(healthJson), nil
}

func GetHealthChecksTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("health_checks",
		mcp.WithDescription("Returns the health checks for a specific service in the Consul cluster."),
		mcp.WithTitleAnnotation("Get health checks for a specific service in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_name",
			mcp.Description("The name of the service to query health checks for."),
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
		mcp.WithString("near",
			mcp.Description("Sort the node list in ascending order based on the estimated round trip time from the given node."),
		),
		mcp.WithString("node_meta",
			mcp.Description("Filter results to only nodes with the specified key/value pairs in their metadata."),
		),
		mcp.WithString("filter",
			mcp.Description("Specifies the expression used to filter the queries results prior to returning the data."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getHealthChecksHandler(ctx, request, logger)
		},
	}
}

func getHealthChecksHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceName, err := request.RequireString("service_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")
	near := request.GetString("near", "")
	nodeMeta := request.GetString("node_meta", "")
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

	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if near != "" {
		queryParams.Set("near", near)
	}
	if nodeMeta != "" {
		queryParams.Set("node-meta", nodeMeta)
	}
	if filter != "" {
		queryParams.Set("filter", filter)
	}

	checksResp, err := consulClient.Get(fmt.Sprintf("health/checks/%s", serviceName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching health checks for service '%s' from consul", serviceName), err)
	}

	// convert checksResp i.e. bytes[] to text
	checksJson := strings.TrimSpace(string(checksResp))
	return mcp.NewToolResultText(checksJson), nil
}

func GetHealthServiceTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("health_service",
		mcp.WithDescription("Returns the health information for all instances of a specific service in the Consul cluster."),
		mcp.WithTitleAnnotation("Get health information for all instances of a specific service in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_name",
			mcp.Description("The name of the service to query health information for."),
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
		mcp.WithString("near",
			mcp.Description("Sort the node list in ascending order based on the estimated round trip time from the given node."),
		),
		mcp.WithString("tag",
			mcp.Description("Filter results to only nodes providing the service with this tag."),
		),
		mcp.WithString("node_meta",
			mcp.Description("Filter results to only nodes with the specified key/value pairs in their metadata."),
		),
		mcp.WithString("passing",
			mcp.Description("Filter results to only return instances with passing health checks. Set to 'true' to enable."),
		),
		mcp.WithString("filter",
			mcp.Description("Specifies the expression used to filter the queries results prior to returning the data."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getHealthServiceHandler(ctx, request, logger)
		},
	}
}

func getHealthServiceHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceName, err := request.RequireString("service_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")
	near := request.GetString("near", "")
	tag := request.GetString("tag", "")
	nodeMeta := request.GetString("node_meta", "")
	passing := request.GetString("passing", "")
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

	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if near != "" {
		queryParams.Set("near", near)
	}
	if tag != "" {
		queryParams.Set("tag", tag)
	}
	if nodeMeta != "" {
		queryParams.Set("node-meta", nodeMeta)
	}
	if passing == "true" {
		queryParams.Set("passing", "true")
	}
	if filter != "" {
		queryParams.Set("filter", filter)
	}

	serviceResp, err := consulClient.Get(fmt.Sprintf("health/service/%s", serviceName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching health information for service '%s' from consul", serviceName), err)
	}

	// convert serviceResp i.e. bytes[] to text
	serviceJson := strings.TrimSpace(string(serviceResp))
	return mcp.NewToolResultText(serviceJson), nil
}

func GetHealthConnectTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("health_connect",
		mcp.WithDescription("Returns the health information for Connect-enabled service instances in the Consul cluster."),
		mcp.WithTitleAnnotation("Get health information for Connect-enabled service instances in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_name",
			mcp.Description("The name of the Connect-enabled service to query health information for."),
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
		mcp.WithString("near",
			mcp.Description("Sort the node list in ascending order based on the estimated round trip time from the given node."),
		),
		mcp.WithString("tag",
			mcp.Description("Filter results to only nodes providing the service with this tag."),
		),
		mcp.WithString("node_meta",
			mcp.Description("Filter results to only nodes with the specified key/value pairs in their metadata."),
		),
		mcp.WithString("passing",
			mcp.Description("Filter results to only return instances with passing health checks. Set to 'true' to enable."),
		),
		mcp.WithString("filter",
			mcp.Description("Specifies the expression used to filter the queries results prior to returning the data."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getHealthConnectHandler(ctx, request, logger)
		},
	}
}

func getHealthConnectHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceName, err := request.RequireString("service_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")
	near := request.GetString("near", "")
	tag := request.GetString("tag", "")
	nodeMeta := request.GetString("node_meta", "")
	passing := request.GetString("passing", "")
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

	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if near != "" {
		queryParams.Set("near", near)
	}
	if tag != "" {
		queryParams.Set("tag", tag)
	}
	if nodeMeta != "" {
		queryParams.Set("node-meta", nodeMeta)
	}
	if passing == "true" {
		queryParams.Set("passing", "true")
	}
	if filter != "" {
		queryParams.Set("filter", filter)
	}

	connectResp, err := consulClient.Get(fmt.Sprintf("health/connect/%s", serviceName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching health information for Connect service '%s' from consul", serviceName), err)
	}

	// convert connectResp i.e. bytes[] to text
	connectJson := strings.TrimSpace(string(connectResp))
	return mcp.NewToolResultText(connectJson), nil
}

func GetHealthIngressTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("health_ingress",
		mcp.WithDescription("Returns the health information for ingress gateway instances in the Consul cluster."),
		mcp.WithTitleAnnotation("Get health information for ingress gateway instances in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service_name",
			mcp.Description("The name of the ingress gateway service to query health information for."),
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
		mcp.WithString("near",
			mcp.Description("Sort the node list in ascending order based on the estimated round trip time from the given node."),
		),
		mcp.WithString("tag",
			mcp.Description("Filter results to only nodes providing the service with this tag."),
		),
		mcp.WithString("node_meta",
			mcp.Description("Filter results to only nodes with the specified key/value pairs in their metadata."),
		),
		mcp.WithString("passing",
			mcp.Description("Filter results to only return instances with passing health checks. Set to 'true' to enable."),
		),
		mcp.WithString("filter",
			mcp.Description("Specifies the expression used to filter the queries results prior to returning the data."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getHealthIngressHandler(ctx, request, logger)
		},
	}
}

func getHealthIngressHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	serviceName, err := request.RequireString("service_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")
	near := request.GetString("near", "")
	tag := request.GetString("tag", "")
	nodeMeta := request.GetString("node_meta", "")
	passing := request.GetString("passing", "")
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

	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if near != "" {
		queryParams.Set("near", near)
	}
	if tag != "" {
		queryParams.Set("tag", tag)
	}
	if nodeMeta != "" {
		queryParams.Set("node-meta", nodeMeta)
	}
	if passing == "true" {
		queryParams.Set("passing", "true")
	}
	if filter != "" {
		queryParams.Set("filter", filter)
	}

	ingressResp, err := consulClient.Get(fmt.Sprintf("health/ingress/%s", serviceName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching health information for ingress gateway '%s' from consul", serviceName), err)
	}

	// convert ingressResp i.e. bytes[] to text
	ingressJson := strings.TrimSpace(string(ingressResp))
	return mcp.NewToolResultText(ingressJson), nil
}

func GetHealthStateTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("health_state",
		mcp.WithDescription("Returns the health checks in a specific state across all services and nodes in the Consul cluster."),
		mcp.WithTitleAnnotation("Get health checks in a specific state across all services and nodes in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("state",
			mcp.Description("The state of health checks to filter by (any, passing, warning, critical)."),
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
		mcp.WithString("near",
			mcp.Description("Sort the node list in ascending order based on the estimated round trip time from the given node."),
		),
		mcp.WithString("node_meta",
			mcp.Description("Filter results to only nodes with the specified key/value pairs in their metadata."),
		),
		mcp.WithString("filter",
			mcp.Description("Specifies the expression used to filter the queries results prior to returning the data."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getHealthStateHandler(ctx, request, logger)
		},
	}
}

func getHealthStateHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	state, err := request.RequireString("state")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: state is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Validate state parameter
	validStates := []string{"any", "passing", "warning", "critical"}
	validState := false
	for _, validStateValue := range validStates {
		if state == validStateValue {
			validState = true
			break
		}
	}
	if !validState {
		return nil, utils.LogAndReturnError(logger, "invalid 'state' parameter: must be one of 'any', 'passing', 'warning', 'critical'", fmt.Errorf("invalid state parameter: %s", state))
	}

	// Get optional parameters
	dc := request.GetString("dc", "")
	near := request.GetString("near", "")
	nodeMeta := request.GetString("node_meta", "")
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

	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if near != "" {
		queryParams.Set("near", near)
	}
	if nodeMeta != "" {
		queryParams.Set("node-meta", nodeMeta)
	}
	if filter != "" {
		queryParams.Set("filter", filter)
	}

	stateResp, err := consulClient.Get(fmt.Sprintf("health/state/%s", state), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching health state '%s' from consul", state), err)
	}

	// convert stateResp i.e. bytes[] to text
	stateJson := strings.TrimSpace(string(stateResp))
	return mcp.NewToolResultText(stateJson), nil
}
