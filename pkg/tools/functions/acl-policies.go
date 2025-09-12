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

func GetACLPolicesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_policies",
		mcp.WithDescription("Returns the list of ACL policies in the Consul cluster."),
		mcp.WithTitleAnnotation("List ACL policies in the Consul cluster"),
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
			return getACLPoliciesHandler(ctx, request, logger)
		},
	}
}

func getACLPoliciesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	uri := (&url.URL{
		Path: "acl/policies",
		RawQuery: url.Values{
			"partition": {ap},
			"ns":        {ns},
		}.Encode(),
	}).String()

	policiesResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching ACL policies list from consul", err)
	}

	// convert policiesResp i.e. bytes[] to text
	policiesJson := strings.TrimSpace(string(policiesResp))
	return mcp.NewToolResultText(policiesJson), nil
}

func GetACLPolicyTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_policy",
		mcp.WithDescription("Returns the details of a specific ACL policy by ID."),
		mcp.WithTitleAnnotation("Get specific ACL policy details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("policy_id",
			mcp.Description("The UUID of the ACL policy to query."),
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
			return getACLPolicyHandler(ctx, request, logger)
		},
	}
}

func getACLPolicyHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	policyId, err := request.RequireString("policy_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: policy_id is required", err)
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
		Path: fmt.Sprintf("acl/policy/%s", policyId),
		RawQuery: url.Values{
			"partition": {ap},
			"ns":        {ns},
		}.Encode(),
	}).String()

	policyResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching ACL policy '%s' details from consul", policyId), err)
	}

	// convert policyResp i.e. bytes[] to text
	policyJson := strings.TrimSpace(string(policyResp))
	return mcp.NewToolResultText(policyJson), nil
}
