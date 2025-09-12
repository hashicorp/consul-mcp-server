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

func GetOperatorLicenseTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("operator_license",
		mcp.WithDescription("Returns the current license information for the Consul cluster."),
		mcp.WithTitleAnnotation("Get license information from the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getOperatorLicenseHandler(ctx, request, logger)
		},
	}
}

func getOperatorLicenseHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	// Get optional parameters
	dc := request.GetString("dc", "")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{}
	if dc != "" {
		queryParams.Set("dc", dc)
	}

	uri := (&url.URL{
		Path:     "operator/license",
		RawQuery: queryParams.Encode(),
	}).String()

	licenseResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching license information from consul operator", err)
	}

	// convert licenseResp i.e. bytes[] to text
	licenseJson := strings.TrimSpace(string(licenseResp))
	return mcp.NewToolResultText(licenseJson), nil
}
