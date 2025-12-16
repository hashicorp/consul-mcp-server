// Copyright IBM Corp. 2025
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/consul-mcp-server/pkg/client"
	"github.com/hashicorp/consul-mcp-server/pkg/utils"
)

func GetIdentity(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("cluster_identity",
		mcp.WithDescription("Returns the Cluster's identity information."),
		mcp.WithTitleAnnotation("Get Consul cluster identity information"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("pem",
			mcp.Description("Set to 'true' to get the response in PEM format."),
			mcp.DefaultString("false"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getIdentityHandler(ctx, request, logger)
		},
	}
}

func getIdentityHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	pem := request.GetString("pem", "false")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{}
	if pem == "true" {
		queryParams.Set("pem", "true")
	}

	rootsResp, err := consulClient.Get("connect/ca/roots", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching Connect CA roots from consul", err)
	}

	// Create types.ConnectCARoots from rootsResp
	var caRoots api.CARootList
	if err := json.Unmarshal(rootsResp, &caRoots); err != nil {
		return nil, utils.LogAndReturnError(logger, "failed to unmarshal Connect CA roots response", err)
	}

	// convert rootsResp i.e. bytes[] to text
	rootsText := strings.TrimSpace(caRoots.TrustDomain)
	return mcp.NewToolResultText(rootsText), nil
}
