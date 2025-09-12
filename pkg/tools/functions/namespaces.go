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

func GetNamespacesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("namespaces",
		mcp.WithDescription("Returns the list of namespaces in the Consul cluster."),
		mcp.WithTitleAnnotation("List namespaces in the Consul cluster"),
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
			return getNamespacesHandler(ctx, request, logger)
		},
	}
}

func getNamespacesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	ap := request.GetString("admin_partition", "default")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	uri := (&url.URL{
		Path: "namespaces",
		RawQuery: url.Values{
			"partition": {ap},
		}.Encode(),
	}).String()

	namespacesResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching namespaces list from consul", err)
	}

	// convert namespacesResp i.e. bytes[] to text
	namespacesJson := strings.TrimSpace(string(namespacesResp))
	return mcp.NewToolResultText(namespacesJson), nil
}

func GetNamespaceTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("namespace",
		mcp.WithDescription("Returns the details of a specific namespace in the Consul cluster."),
		mcp.WithTitleAnnotation("Get specific namespace details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("namespace_name",
			mcp.Description("The name of the namespace to query."),
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
			return getNamespaceHandler(ctx, request, logger)
		},
	}
}

func getNamespaceHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	namespaceName, err := request.RequireString("namespace_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: namespace_name is required", err)
	}

	ap := request.GetString("admin_partition", "default")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	uri := (&url.URL{
		Path: fmt.Sprintf("namespace/%s", namespaceName),
		RawQuery: url.Values{
			"partition": {ap},
		}.Encode(),
	}).String()

	namespaceResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching namespace '%s' details from consul", namespaceName), err)
	}

	// convert namespaceResp i.e. bytes[] to text
	namespaceJson := strings.TrimSpace(string(namespaceResp))
	return mcp.NewToolResultText(namespaceJson), nil
}
