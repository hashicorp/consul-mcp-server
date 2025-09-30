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

func GetOperatorKeyringTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("operator_keyring",
		mcp.WithDescription("Returns the gossip encryption keyring information for the Consul cluster."),
		mcp.WithTitleAnnotation("Get gossip encryption keyring information from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("relay_factor",
			mcp.Description("Setting this to a non-zero value will cause nodes to relay their response through this many randomly-chosen other nodes in the cluster."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getOperatorKeyringHandler(ctx, request, logger)
		},
	}
}

func getOperatorKeyringHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get optional parameters
	relayFactor := request.GetString("relay_factor", "")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{}
	if relayFactor != "" {
		queryParams.Set("relay-factor", relayFactor)
	}

	keyringResp, err := consulClient.Get("operator/keyring", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching keyring information from consul operator", err)
	}

	// convert keyringResp i.e. bytes[] to text
	keyringJson := strings.TrimSpace(string(keyringResp))
	return mcp.NewToolResultText(keyringJson), nil
}
