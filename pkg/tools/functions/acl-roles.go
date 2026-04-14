// Copyright (c) HashiCorp, Inc.
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

func GetACLRolesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_roles",
		mcp.WithDescription("Returns the list of ACL roles in the Consul cluster."),
		mcp.WithTitleAnnotation("List ACL roles in the Consul cluster"),
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
			return getACLRolesHandler(ctx, request, logger)
		},
	}
}

func getACLRolesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
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

	rolesResp, err := consulClient.Get("acl/roles", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching ACL roles list from consul", err)
	}

	// convert rolesResp i.e. bytes[] to text
	rolesJson := strings.TrimSpace(string(rolesResp))
	return mcp.NewToolResultText(rolesJson), nil
}

func GetACLRoleTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_role",
		mcp.WithDescription("Returns the details of a specific ACL role by ID."),
		mcp.WithTitleAnnotation("Get specific ACL role details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("role_id",
			mcp.Description("The UUID of the ACL role to query."),
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
			return getACLRoleHandler(ctx, request, logger)
		},
	}
}

func getACLRoleHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	roleId, err := request.RequireString("role_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: role_id is required", err)
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

	roleResp, err := consulClient.Get(fmt.Sprintf("acl/role/%s", roleId), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching ACL role '%s' details from consul", roleId), err)
	}

	// convert roleResp i.e. bytes[] to text
	roleJson := strings.TrimSpace(string(roleResp))
	return mcp.NewToolResultText(roleJson), nil
}
