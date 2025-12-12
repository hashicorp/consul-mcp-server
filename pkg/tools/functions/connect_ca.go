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

func GetConnectCARootsTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("connect_ca_roots",
		mcp.WithDescription("Returns the trusted CA root certificates for Connect in the Consul cluster."),
		mcp.WithTitleAnnotation("Get Connect CA root certificates from the Consul cluster"),
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
			return getConnectCARootsHandler(ctx, request, logger)
		},
	}
}

func getConnectCARootsHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
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

	// convert rootsResp i.e. bytes[] to text
	rootsJson := strings.TrimSpace(string(rootsResp))
	return mcp.NewToolResultText(rootsJson), nil
}

func GetConnectCAConfigurationTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("connect_ca_configuration",
		mcp.WithDescription("Returns the Connect CA configuration for the Consul cluster."),
		mcp.WithTitleAnnotation("Get Connect CA configuration from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getConnectCAConfigurationHandler(ctx, request, logger)
		},
	}
}

func getConnectCAConfigurationHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	configResp, err := consulClient.Get("connect/ca/configuration", nil)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching Connect CA configuration from consul", err)
	}

	// convert configResp i.e. bytes[] to text
	configJson := strings.TrimSpace(string(configResp))
	return mcp.NewToolResultText(configJson), nil
}
