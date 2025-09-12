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

func GetAgentConnectCATool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_connect_ca",
		mcp.WithDescription("Returns the CA certificate bundle from the Connect CA that can be used to verify a TLS connection with the local agent."),
		mcp.WithTitleAnnotation("Get Connect CA certificate bundle"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentConnectCAHandler(ctx, request, logger)
		},
	}
}

func getAgentConnectCAHandler(ctx context.Context, _ mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	uri := "agent/connect/ca"

	caResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "fetching Connect CA certificate bundle", err)
	}

	caOutput := strings.TrimSpace(string(caResp))
	return mcp.NewToolResultText(caOutput), nil
}

func GetAgentConnectAuthorizeTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_connect_authorize",
		mcp.WithDescription("Tests whether a connection is authorized between two services based on Connect intentions."),
		mcp.WithTitleAnnotation("Test Connect authorization between services"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("target",
			mcp.Description("The target service name to check authorization for."),
			mcp.Required(),
		),
		mcp.WithString("client_cert_uri",
			mcp.Description("The client certificate URI to check authorization for."),
			mcp.Required(),
		),
		mcp.WithString("client_cert_serial",
			mcp.Description("The client certificate serial number."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentConnectAuthorizeHandler(ctx, request, logger)
		},
	}
}

func getAgentConnectAuthorizeHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	target, err := request.RequireString("target")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: target is required", err)
	}

	clientCertURI, err := request.RequireString("client_cert_uri")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: client_cert_uri is required", err)
	}

	clientCertSerial := request.GetString("client_cert_serial", "")

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{
		"target":          {target},
		"client_cert_uri": {clientCertURI},
	}

	if clientCertSerial != "" {
		queryParams.Set("client_cert_serial", clientCertSerial)
	}

	uri := (&url.URL{
		Path:     "agent/connect/authorize",
		RawQuery: queryParams.Encode(),
	}).String()

	authResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "checking Connect authorization", err)
	}

	authJson := strings.TrimSpace(string(authResp))
	return mcp.NewToolResultText(authJson), nil
}

func GetAgentConnectProxyConfigTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_connect_proxy_config",
		mcp.WithDescription("Returns the configuration for a Connect proxy."),
		mcp.WithTitleAnnotation("Get Connect proxy configuration"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("proxy_service_id",
			mcp.Description("The ID of the proxy service to get configuration for."),
			mcp.Required(),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentConnectProxyConfigHandler(ctx, request, logger)
		},
	}
}

func getAgentConnectProxyConfigHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	proxyServiceID, err := request.RequireString("proxy_service_id")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: proxy_service_id is required", err)
	}

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	uri := fmt.Sprintf("agent/connect/proxy/%s", proxyServiceID)

	proxyResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching Connect proxy configuration for '%s'", proxyServiceID), err)
	}

	proxyJson := strings.TrimSpace(string(proxyResp))
	return mcp.NewToolResultText(proxyJson), nil
}

func GetAgentConnectLeafCertTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("agent_connect_leaf_cert",
		mcp.WithDescription("Returns the leaf certificate representing the specified service."),
		mcp.WithTitleAnnotation("Get Connect leaf certificate for service"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("service",
			mcp.Description("The name of the service to get the leaf certificate for."),
			mcp.Required(),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getAgentConnectLeafCertHandler(ctx, request, logger)
		},
	}
}

func getAgentConnectLeafCertHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	service, err := request.RequireString("service")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: service is required", err)
	}

	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	queryParams := url.Values{
		"service": {service},
	}

	uri := (&url.URL{
		Path:     "agent/connect/leaf",
		RawQuery: queryParams.Encode(),
	}).String()

	leafResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching Connect leaf certificate for service '%s'", service), err)
	}

	leafJson := strings.TrimSpace(string(leafResp))
	return mcp.NewToolResultText(leafJson), nil
}
