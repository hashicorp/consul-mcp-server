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

func GetKVTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("kv_get",
		mcp.WithDescription("Returns the value for a specific key from the Consul KV store."),
		mcp.WithTitleAnnotation("Get a specific key from the Consul KV store"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("key",
			mcp.Description("The key to retrieve from the KV store."),
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
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
		mcp.WithString("raw",
			mcp.Description("Return raw value without JSON metadata. Set to 'true' to enable."),
			mcp.DefaultString("false"),
		),
		mcp.WithString("separator",
			mcp.Description("Character to use as list separator when listing keys."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getKVHandler(ctx, request, logger)
		},
	}
}

func getKVHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	key, err := request.RequireString("key")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: key is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")
	raw := request.GetString("raw", "false")
	separator := request.GetString("separator", "")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{
		"partition": {ap},
		"ns":        {ns},
	}

	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if raw == "true" {
		queryParams.Set("raw", "true")
	}
	if separator != "" {
		queryParams.Set("separator", separator)
	}

	uri := (&url.URL{
		Path:     fmt.Sprintf("kv/%s", key),
		RawQuery: queryParams.Encode(),
	}).String()

	kvResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching key '%s' from consul KV store", key), err)
	}

	// convert kvResp i.e. bytes[] to text
	kvText := strings.TrimSpace(string(kvResp))
	return mcp.NewToolResultText(kvText), nil
}

func GetKVKeysTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("kv_keys",
		mcp.WithDescription("Returns the list of keys with a given prefix from the Consul KV store."),
		mcp.WithTitleAnnotation("List keys with a prefix from the Consul KV store"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("prefix",
			mcp.Description("The key prefix to list keys for. Use empty string to list all keys."),
			mcp.DefaultString(""),
		),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("namespace",
			mcp.Description("The consul namespace to query."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
		mcp.WithString("separator",
			mcp.Description("Character to use as list separator when listing keys."),
			mcp.DefaultString("/"),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getKVKeysHandler(ctx, request, logger)
		},
	}
}

func getKVKeysHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	prefix := request.GetString("prefix", "")



	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")
	separator := request.GetString("separator", "/")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{
		"keys":      {"true"},
		"partition": {ap},
		"ns":        {ns},
	}

	if dc != "" {
		queryParams.Set("dc", dc)
	}
	if separator != "" {
		queryParams.Set("separator", separator)
	}

	var uri string
	if prefix == "" {
		uri = (&url.URL{
			Path:     "kv/",
			RawQuery: queryParams.Encode(),
		}).String()
	} else {
		uri = (&url.URL{
			Path:     fmt.Sprintf("kv/%s", prefix),
			RawQuery: queryParams.Encode(),
		}).String()
	}

	keysResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching keys with prefix '%s' from consul KV store", prefix), err)
	}

	// convert keysResp i.e. bytes[] to text
	keysJson := strings.TrimSpace(string(keysResp))
	return mcp.NewToolResultText(keysJson), nil
}

func GetKVRecursiveTool(logger *log.Logger) server.ServerTool {
	tool := mcp.NewTool("kv_recursive",
		mcp.WithDescription("Returns all key-value pairs under a given prefix recursively from the Consul KV store."),
		mcp.WithTitleAnnotation("Get all key-value pairs under a prefix recursively from the Consul KV store"),
		mcp.WithOpenWorldHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithString("prefix",
			mcp.Description("The key prefix to retrieve recursively. Use empty string to get all keys."),
			mcp.DefaultString(""),
		),
		mcp.WithString("admin_partition",
			mcp.Description("The consul admin partition to query."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("namespace",
			mcp.Description("The consul namespace to query."),
			mcp.DefaultString("default"),
		),
		mcp.WithString("dc",
			mcp.Description("The datacenter to query. If not provided, the datacenter of the agent is queried."),
		),
	)
	return server.ServerTool{
		Tool: tool,
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return getKVRecursiveHandler(ctx, request, logger)
		},
	}
}

func getKVRecursiveHandler(ctx context.Context, request mcp.CallToolRequest, logger *log.Logger) (*mcp.CallToolResult, error) {
	prefix, err := request.RequireString("prefix")
	if err != nil {
		return nil, utils.LogAndReturnError(logger, "required input: prefix is required", err)
	}

	ap := request.GetString("admin_partition", "default")
	ns := request.GetString("namespace", "default")

	// Get optional parameters
	dc := request.GetString("dc", "")

	// Get a simple http client to access the consul API
	consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get http client for consul API")
		return mcp.NewToolResultError(fmt.Sprintf("failed to get http client for consul API: %v", err)), nil
	}

	// Build query parameters
	queryParams := url.Values{
		"recurse":   {"true"},
		"partition": {ap},
		"ns":        {ns},
	}

	if dc != "" {
		queryParams.Set("dc", dc)
	}

	var uri string
	if prefix == "" {
		uri = (&url.URL{
			Path:     "kv/",
			RawQuery: queryParams.Encode(),
		}).String()
	} else {
		uri = (&url.URL{
			Path:     fmt.Sprintf("kv/%s", prefix),
			RawQuery: queryParams.Encode(),
		}).String()
	}

	recursiveResp, err := consulClient.Get(uri)
	if err != nil {
		return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching keys recursively with prefix '%s' from consul KV store", prefix), err)
	}

	// convert recursiveResp i.e. bytes[] to text
	recursiveJson := strings.TrimSpace(string(recursiveResp))
	return mcp.NewToolResultText(recursiveJson), nil
}
