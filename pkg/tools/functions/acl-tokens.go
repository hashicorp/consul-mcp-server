// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul-mcp-server/pkg/client"
	"github.com/hashicorp/consul-mcp-server/pkg/utils"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetACLTokensTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_tokens",
		mcp.WithDescription("Returns the list of ACL tokens in the Consul cluster."),
		mcp.WithTitleAnnotation("List ACL tokens in the Consul cluster"),
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
		mcp.WithString("policy",
			mcp.Description("Filter results to only tokens that have the specified policy ID."),
		),
		mcp.WithString("role",
			mcp.Description("Filter results to only tokens that have the specified role ID."),
		),
		mcp.WithString("authmethod",
			mcp.Description("Filter results to only tokens that were created by the specified auth method."),
		),
		mcp.WithString("authmethod-ns",
			mcp.Description("The namespace of the auth method to filter by."),
		),
		mcp.WithString("filter",
			mcp.Description("Specifies the expression used to filter the queries results prior to returning the data."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getACLTokensHandler(ctx, request, logger)
		},
	}
}

func getACLTokensHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")
	policy := request.GetString("policy", "")
	role := request.GetString("role", "")
	authMethod := request.GetString("authmethod", "")
	authMethodNs := request.GetString("authmethod-ns", "")
	filter := request.GetString("filter", "")

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

	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if policy != "" {
		queryParams.Set("policy", policy)
	}
	if role != "" {
		queryParams.Set("role", role)
	}
	if authMethod != "" {
		queryParams.Set("authmethod", authMethod)
	}
	if authMethodNs != "" {
		queryParams.Set("authmethod-ns", authMethodNs)
	}
	if filter != "" {
		queryParams.Set("filter", filter)
	}

	tokensResp, err := consulClient.Get("acl/tokens", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching ACL tokens list from consul", err)
	}

	// Parse the JSON response using official Consul SDK types
	var tokens []*api.ACLToken
	if err := json.Unmarshal(tokensResp, &tokens); err != nil {
		return nil, utils.LogAndReturnError(logger, "parsing ACL token self response", err)
	}

	for _, token := range tokens {
		// Redact SecretID field
		if token.SecretID != "" && token.SecretID != "anonymous" {
			token.SecretID = "[REDACTED]"
		}
	}

	// parse back to JSON
	if len(tokens) == 0 {
		tokensResp = []byte("[]")
	} else {
		tokensResp, err = json.MarshalIndent(tokens, "", "  ")
		if err != nil {
			return nil, utils.LogAndReturnError(logger, "marshaling redacted ACL tokens", err)
		}
	}

	// convert tokensResp i.e. bytes[] to text
	tokensJson := strings.TrimSpace(string(tokensResp))
	return mcp.NewToolResultText(tokensJson), nil
}

func GetACLTokenTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_token",
		mcp.WithDescription("Returns the details of a specific ACL token by ID."),
		mcp.WithTitleAnnotation("Get specific ACL token details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("token_id",
			mcp.Description("The UUID of the ACL token to query."),
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
			return getACLTokenHandler(ctx, request, logger)
		},
	}
}

func getACLTokenHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	tokenId, err := request.RequireString("token_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: token_id is required", err)
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

	tokenResp, err := consulClient.Get(fmt.Sprintf("acl/token/%s", tokenId), queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching ACL token '%s' details from consul", tokenId), err)
	}

	// convert tokenResp i.e. bytes[] to text
	tokenJson := strings.TrimSpace(string(tokenResp))
	return mcp.NewToolResultText(tokenJson), nil
}

func GetACLTokenSelfTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("acl_token_self",
		mcp.WithDescription("Returns the details of the current ACL token being used for authentication."),
		mcp.WithTitleAnnotation("Get current ACL token details from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getACLTokenSelfHandler(ctx, request, logger)
		},
	}
}

func getACLTokenSelfHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	uri := "acl/token/self"

	tokenResp, err := consulClient.Get(uri, nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching current ACL token details from consul", err)
	}

	// Parse the JSON response using official Consul SDK types
	var token *api.ACLToken
	if err := json.Unmarshal(tokenResp, &token); err != nil {
		return nil, utils.LogAndReturnError(logger, "parsing ACL token self response", err)
	}

	// Redact SecretID field
	if token.SecretID != "" && token.SecretID != "anonymous" {
		token.SecretID = "[REDACTED]"
	}

	// Marshal back to JSON
	redactedJson, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "marshaling redacted ACL token self", err)
	}

	return mcp.NewToolResultText(string(redactedJson)), nil
}
