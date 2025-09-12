// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul-mcp-server/pkg/utils"
	"github.com/hashicorp/consul-mcp-server/version"
	"github.com/mark3labs/mcp-go/server"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	ConsulSkipTLSVerify = "CONSUL_SKIP_VERIFY"
)

type ConsulHttpClient struct {
	SessionID string
	Address   string
	Token     string
	// Add other Consul-specific fields as needed

	client *http.Client
	Logger *log.Logger
}

var (
	activeHttpClients sync.Map
)

// NewConsulClientFromContext creates a new Consul client from the given context
func NewConsulClientFromContext(ctx context.Context, logger *log.Logger) (*ConsulHttpClient, error) {
	sessionId, ok := ctx.Value("session_id").(string)
	if !ok || sessionId == "" {
		return nil, fmt.Errorf("session_id not found in context")
	}
	return NewConsulClient(ctx, sessionId, logger), nil
}

// DeleteHttpClient removes the HTTP client for the given session
func DeleteConsulHttpClientForSession(sessionId string) {
	activeHttpClients.Delete(sessionId)
}

// NewConsulClient creates a new Consul client for the given session
func NewConsulClient(ctx context.Context, sessionId string, logger *log.Logger) *ConsulHttpClient {
	address := utils.GetEnv("CONSUL_HTTP_ADDR", "http://localhost:8500")
	// Ensure the address does not have a trailing slash
	if address[len(address)-1] == '/' {
		address = address[:len(address)-1]
	}

	token := utils.GetEnv("CONSUL_HTTP_TOKEN", "")

	// override the address and token from session context if available
	if addr, ok := ctx.Value("consul_address").(string); ok && addr != "" {
		address = addr
	}

	if tkn, ok := ctx.Value("consul_token").(string); ok && tkn != "" {
		token = tkn
	}

	httpClient := createHTTPClient(parseSkipTLSVerify(ctx), logger)

	consulClient := &ConsulHttpClient{
		SessionID: sessionId,
		Address:   address,
		Token:     token,
		client:    httpClient,
		Logger:    logger,
	}

	activeHttpClients.Store(sessionId, consulClient)

	return consulClient
}

// GetConsulHttpClient retrieves the Consul client for the given session
func GetConsulHttpClient(sessionId string) *ConsulHttpClient {
	if value, ok := activeHttpClients.Load(sessionId); ok {
		return value.(*ConsulHttpClient)
	}
	return nil
}

// GetHttpClientFromContext extracts HTTP client from the MCP context
func GetGetConsulHttpClientFromContext(ctx context.Context, logger *log.Logger) (*ConsulHttpClient, error) {
	session := server.ClientSessionFromContext(ctx)
	if session == nil {
		return nil, fmt.Errorf("no active session")
	}

	// Try to get existing client
	client := GetConsulHttpClient(session.SessionID())
	if client != nil {
		return client, nil
	}

	logger.Warnf("HTTP client not found, creating a new one")
	return NewConsulClient(ctx, session.SessionID(), logger), nil
}

func (c *ConsulHttpClient) Put(uri string, data interface{}, callOptions ...string) ([]byte, error) {
	return c.call("PUT", uri, data, callOptions...)
}

func (c *ConsulHttpClient) Post(uri string, data interface{}, callOptions ...string) ([]byte, error) {
	return c.call("POST", uri, data, callOptions...)
}

func (c *ConsulHttpClient) call(method string, uri string, data interface{}, callOptions ...string) ([]byte, error) {
	ver := "v1"
	if len(callOptions) > 0 {
		ver = callOptions[0] // API version will be the first optional arg to this function
	}

	parsedURL, err := url.Parse(fmt.Sprintf("%s/%s/%s", c.Address, ver, uri))
	if err != nil {
		return nil, fmt.Errorf("error parsing the URL: %w", err)
	}
	c.Logger.Debugf("Requested URL: %s", parsedURL)

	var reqBody io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error marshalling data to JSON: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
		c.Logger.Tracef("Request body: %s", string(jsonData))
	}

	req, err := http.NewRequest(method, parsedURL.String(), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("consul-mcp-server/%s", version.GetHumanVersion()))
	if c.Token != "" {
		req.Header.Set("X-Consul-Token", c.Token)
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", "404 Not Found")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Logger.Errorf("error closing response body: %v", err)
		}
	}(resp.Body)
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.Logger.Debugf("Response status: %s", resp.Status)
	c.Logger.Tracef("Response body: %s", string(body))
	return body, nil
}

func (c *ConsulHttpClient) Get(uri string, callOptions ...string) ([]byte, error) {
	ver := "v1"
	if len(callOptions) > 0 {
		ver = callOptions[0] // API version will be the first optional arg to this function
	}

	parsedURL, err := url.Parse(fmt.Sprintf("%s/%s/%s", c.Address, ver, uri))
	if err != nil {
		return nil, fmt.Errorf("error parsing the URL: %w", err)
	}
	c.Logger.Debugf("Requested URL: %s", parsedURL)

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("consul-mcp-server/%s", version.GetHumanVersion()))
	if c.Token != "" {
		req.Header.Set("X-Consul-Token", c.Token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", "404 Not Found")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Logger.Errorf("error closing response body: %v", err)
		}
	}(resp.Body)
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.Logger.Debugf("Response status: %s", resp.Status)
	c.Logger.Tracef("Response body: %s", string(body))
	return body, nil
}

func parseSkipTLSVerify(ctx context.Context) bool {
	skipTLSVerifyStr, ok := ctx.Value(contextKey(ConsulSkipTLSVerify)).(string)
	if ok && skipTLSVerifyStr != "" {
		skipTLSVerify, err := strconv.ParseBool(skipTLSVerifyStr)
		if err == nil {
			return skipTLSVerify
		}
	}
	return false
}

// NewHttpClient creates a http.Client with optional TLS verification skip
func NewHttpClientFromContext(ctx context.Context, logger *log.Logger) *http.Client {
	client := createHTTPClient(true, logger)
	sessionId, ok := ctx.Value("session_id").(string)
	if !ok || sessionId == "" {
		logger.WithField("session_id", sessionId).Info("Created HTTP client")
	}
	return client
}
