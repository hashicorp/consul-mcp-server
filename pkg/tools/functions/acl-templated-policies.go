// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package functions

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/consul-mcp-server/pkg/client"
	"github.com/hashicorp/consul-mcp-server/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetACLTemplatedPoliciesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_templated_policies",
		mcp.WithDescription("Returns the list of ACL templated policies in the Consul cluster."),
		mcp.WithTitleAnnotation("List ACL templated policies in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getACLTemplatedPoliciesHandler(ctx, request, logger)
		},
	}
}

func getACLTemplatedPoliciesHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	templatedPoliciesResp, err := consulClient.Get("acl/templated-policies", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching ACL templated policies list from consul", err)
	}

	// convert templatedPoliciesResp i.e. bytes[] to text
	templatedPoliciesJson := strings.TrimSpace(string(templatedPoliciesResp))
	return mcp.NewToolResultText(templatedPoliciesJson), nil
}

func GetACLTemplatedPolicyTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_templated_policy",
		mcp.WithDescription("Returns the details of a specific ACL templated policy by name."),
		mcp.WithTitleAnnotation("Get specific ACL templated policy details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("policy_name",
			mcp.Description("The name of the ACL templated policy to query."),
			mcp.Required(),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getACLTemplatedPolicyHandler(ctx, request, logger)
		},
	}
}

func getACLTemplatedPolicyHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	policyName, err := request.RequireString("policy_name")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: policy_name is required", err)
	}

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	templatedPolicyResp, err := consulClient.Get(fmt.Sprintf("acl/templated-policy/%s", policyName), nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching ACL templated policy '%s' details from consul", policyName), err)
	}

	// convert templatedPolicyResp i.e. bytes[] to text
	templatedPolicyJson := strings.TrimSpace(string(templatedPolicyResp))
	return mcp.NewToolResultText(templatedPolicyJson), nil
}
