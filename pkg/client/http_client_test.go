// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateHTTPClient(t *testing.T) {
	t.Run("with TLS verification", func(t *testing.T) {
		client := createHTTPClient(false, log.New())
		assert.NotNil(t, client)

		// Check that TLS config is properly set
		roundTripper, ok := client.Transport.(*retryablehttp.RoundTripper)
		transport := roundTripper.Client.HTTPClient.Transport.(*http.Transport)
		assert.True(t, ok)
		assert.NotNil(t, transport)
		assert.NotNil(t, transport.TLSClientConfig)
		assert.False(t, transport.TLSClientConfig.InsecureSkipVerify)
	})

	t.Run("with TLS skip verify", func(t *testing.T) {
		client := createHTTPClient(true, log.New())
		assert.NotNil(t, client)

		// Check that TLS config skips verification
		roundTripper, ok := client.Transport.(*retryablehttp.RoundTripper)
		transport := roundTripper.Client.HTTPClient.Transport.(*http.Transport)
		assert.True(t, ok)
		assert.NotNil(t, transport.TLSClientConfig)
		assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
	})
}

func TestHTTPClientWithCustomTransport(t *testing.T) {
	t.Run("custom timeout", func(t *testing.T) {
		client := createHTTPClient(false, log.New())

		roundTripper, ok := client.Transport.(*retryablehttp.RoundTripper)
		assert.True(t, ok)
		assert.Equal(t, 10*time.Second, roundTripper.Client.HTTPClient.Timeout)
	})

	t.Run("custom transport settings", func(t *testing.T) {
		client := createHTTPClient(false, log.New())
		roundTripper, ok := client.Transport.(*retryablehttp.RoundTripper)
		assert.True(t, ok)

		// Access the retryable client through the RoundTripper
		retryableClient := roundTripper.Client
		assert.Equal(t, 3, retryableClient.RetryMax)

		// Access the underlying HTTP transport
		transport, ok := retryableClient.HTTPClient.Transport.(*http.Transport)
		assert.True(t, ok)
		assert.NotNil(t, transport.TLSClientConfig)
	})
}

func TestNewHttpClientFromContext(t *testing.T) {
	t.Run("with session_id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "session_id", "test-session")
		client := NewHttpClientFromContext(ctx, log.New())
		assert.NotNil(t, client)
		assert.IsType(t, &http.Client{}, client)
	})

	t.Run("without session_id", func(t *testing.T) {
		ctx := context.Background()
		client := NewHttpClientFromContext(ctx, log.New())
		assert.NotNil(t, client)
		assert.IsType(t, &http.Client{}, client)
	})
}

func TestContextKey(t *testing.T) {
	key1 := contextKey("test")
	key2 := contextKey("test")
	key3 := contextKey("other")

	// Same string should create equal keys
	assert.Equal(t, key1, key2)

	// Different strings should create different keys
	assert.NotEqual(t, key1, key3)

	// Test string representation
	assert.Equal(t, "test", string(key1))
}

func TestHTTPClientConfiguration(t *testing.T) {
	// Test with custom TLS config
	client := createHTTPClient(true, log.New())

	roundTripper, ok := client.Transport.(*retryablehttp.RoundTripper)
	assert.True(t, ok)

	transport, ok := roundTripper.Client.HTTPClient.Transport.(*http.Transport)
	assert.True(t, ok)

	// Verify transport configuration
	assert.NotNil(t, transport.TLSClientConfig)
	assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)

	// Test timeout configuration
	assert.Equal(t, 10*time.Second, roundTripper.Client.HTTPClient.Timeout)
}

func TestTLSConfigDefaults(t *testing.T) {
	client := createHTTPClient(false, log.New())
	roundTripper, ok := client.Transport.(*retryablehttp.RoundTripper)
	assert.True(t, ok)

	transport, ok := roundTripper.Client.HTTPClient.Transport.(*http.Transport)
	assert.True(t, ok)
	tlsConfig := transport.TLSClientConfig

	// Test default TLS configuration
	assert.NotNil(t, tlsConfig)
	assert.False(t, tlsConfig.InsecureSkipVerify)
}
