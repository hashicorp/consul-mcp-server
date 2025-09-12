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

func GetACLBindingRulesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_binding_rules",
		mcp.WithDescription("Returns the list of ACL binding rules in the Consul cluster."),
		mcp.WithTitleAnnotation("List ACL binding rules in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("auth_method",
			mcp.Description("Filter binding rules by auth method name."),
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
			return getACLBindingRulesHandler(ctx, request, logger)
		},
	}
}

func getACLBindingRulesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	authMethod := request.GetString("auth_method", "")

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

	if authMethod != "" {
		queryParams.Set("authmethod", authMethod)
	}

	uri := (&url.URL{
		Path:     "acl/binding-rules",
		RawQuery: queryParams.Encode(),
	}).String()

	bindingRulesResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching ACL binding rules list from consul", err)
	}

	// convert bindingRulesResp i.e. bytes[] to text
	bindingRulesJson := strings.TrimSpace(string(bindingRulesResp))
	return mcp.NewToolResultText(bindingRulesJson), nil
}

func GetACLBindingRuleTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_binding_rule",
		mcp.WithDescription("Returns the details of a specific ACL binding rule by ID."),
		mcp.WithTitleAnnotation("Get specific ACL binding rule details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("binding_rule_id",
			mcp.Description("The UUID of the ACL binding rule to query."),
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
			return getACLBindingRuleHandler(ctx, request, logger)
		},
	}
}

func getACLBindingRuleHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	bindingRuleId, err := request.RequireString("binding_rule_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: binding_rule_id is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	uri := (&url.URL{
		Path: fmt.Sprintf("acl/binding-rule/%s", bindingRuleId),
		RawQuery: url.Values{
			"partition": {ap},
			"ns":        {ns},
		}.Encode(),
	}).String()

	bindingRuleResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching ACL binding rule '%s' details from consul", bindingRuleId), err)
	}

	// convert bindingRuleResp i.e. bytes[] to text
	bindingRuleJson := strings.TrimSpace(string(bindingRuleResp))
	return mcp.NewToolResultText(bindingRuleJson), nil
}
