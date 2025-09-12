// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	uri := (&url.URL{
		Path: "acl/tokens",
		RawQuery: url.Values{
			"partition": {ap},
			"ns":        {ns},
		}.Encode(),
	}).String()

	tokensResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching ACL tokens list from consul", err)
	}

	// Parse the JSON response using official Consul SDK types
	var tokens []*api.ACLTokenListEntry
	if err := json.Unmarshal(tokensResp, &tokens); err != nil {
		return nil, utils.LogAndReturnError(logger, "parsing ACL tokens response", err)
	}

	// Redact SecretID fields
	for i := range tokens {
		if tokens[i].SecretID != "" && tokens[i].SecretID != "anonymous" {
			tokens[i].SecretID = "[REDACTED]"
		}
	}

	// Marshal back to JSON
	redactedJson, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "marshaling redacted ACL tokens", err)
	}

	return mcp.NewToolResultText(string(redactedJson)), nil
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

	uri := (&url.URL{
		Path: fmt.Sprintf("acl/token/%s", tokenId),
		RawQuery: url.Values{
			"partition": {ap},
			"ns":        {ns},
		}.Encode(),
	}).String()

	tokenResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching ACL token '%s' details from consul", tokenId), err)
	}

	// Parse the JSON response using official Consul SDK types
	var token *api.ACLToken
	if err := json.Unmarshal(tokenResp, &token); err != nil {
		return nil, utils.LogAndReturnError(logger, "parsing ACL token response", err)
	}

	// Redact SecretID field
	if token.SecretID != "" && token.SecretID != "anonymous" {
		token.SecretID = "[REDACTED]"
	}

	// Marshal back to JSON
	redactedJson, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "marshaling redacted ACL token", err)
	}

	return mcp.NewToolResultText(string(redactedJson)), nil
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

	tokenResp, err := consulClient.Get(uri)
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
