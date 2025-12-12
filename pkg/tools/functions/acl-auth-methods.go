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

func GetACLAuthMethodsTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_auth_methods",
		mcp.WithDescription("Returns the list of ACL auth methods in the Consul cluster."),
		mcp.WithTitleAnnotation("List ACL auth methods in the Consul cluster"),
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
			return getACLAuthMethodsHandler(ctx, request, logger)
		},
	}
}

func getACLAuthMethodsHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
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

	authMethodsResp, err := consulClient.Get("acl/auth-methods", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching ACL auth methods list from consul", err)
	}

	// convert authMethodsResp i.e. bytes[] to text
	authMethodsJson := strings.TrimSpace(string(authMethodsResp))
	return mcp.NewToolResultText(authMethodsJson), nil
}

func GetACLAuthMethodTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_auth_method",
		mcp.WithDescription("Returns the details of a specific ACL auth method by name."),
		mcp.WithTitleAnnotation("Get specific ACL auth method details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("auth_method_name",
			mcp.Description("The name of the ACL auth method to query."),
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
			return getACLAuthMethodHandler(ctx, request, logger)
		},
	}
}

func getACLAuthMethodHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	authMethodName, err := request.RequireString("auth_method_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: auth_method_name is required", err)
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

	authMethodResp, err := consulClient.Get(fmt.Sprintf("acl/auth-method/%s", authMethodName), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching ACL auth method '%s' details from consul", authMethodName), err)
	}

	// convert authMethodResp i.e. bytes[] to text
	authMethodJson := strings.TrimSpace(string(authMethodResp))
	return mcp.NewToolResultText(authMethodJson), nil
}
