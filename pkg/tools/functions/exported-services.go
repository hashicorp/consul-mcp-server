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

func GetExportedServicesTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("exported_services",
		mcp.WithDescription("Returns the list of exported services in the Consul cluster."),
		mcp.WithTitleAnnotation("List exported services in the Consul cluster"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getExportedServicesHandler(ctx, request, logger)
		},
	}
}

func getExportedServicesHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	ap := request.GetString("admin_partition", "default")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{
		"partition": {ap},
	}

	exportedResp, err := consulClient.Get("exported-services", queryParams)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching exported services list from consul", err)
	}

	// convert exportedResp i.e. bytes[] to text
	exportedJson := strings.TrimSpace(string(exportedResp))
	return mcp.NewToolResultText(exportedJson), nil
}
