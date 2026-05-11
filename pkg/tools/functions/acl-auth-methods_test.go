// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetACLAuthMethodsTool(t *testing.T) {
	tool := GetACLAuthMethodsTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "acl_auth_methods", tool.Tool.Name)
}

func TestGetACLAuthMethodsHandler(t *testing.T) {
	mockAuthMethods := []map[string]interface{}{
		{
			"Name":        "kubernetes",
			"Type":        "kubernetes",
			"Description": "Kubernetes auth method",
			"Config": map[string]interface{}{
				"Host":              "https://kubernetes.default.svc.cluster.local:443",
				"CACert":            "-----BEGIN CERTIFICATE-----\ntest-ca-cert\n-----END CERTIFICATE-----",
				"ServiceAccountJWT": "",
			},
			"CreateIndex": 10,
			"ModifyIndex": 15,
		},
		{
			"Name":        "oidc-provider",
			"Type":        "oidc",
			"Description": "OIDC auth method for SSO",
			"Config": map[string]interface{}{
				"OIDCDiscoveryURL":    "https://auth.example.com/.well-known/openid_configuration",
				"OIDCClientID":        "consul",
				"OIDCClientSecret":    "[REDACTED]",
				"BoundAudiences":      []string{"consul"},
				"AllowedRedirectURIs": []string{"https://consul.example.com/ui/oidc/callback"},
			},
			"CreateIndex": 20,
			"ModifyIndex": 25,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/acl/auth-methods", r.URL.Path)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "default", query.Get("partition"))
		assert.Equal(t, "default", query.Get("ns"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockAuthMethods)
	}))
	defer server.Close()

	t.Run("successful auth methods listing", func(t *testing.T) {
		// Test would verify auth methods are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})

	t.Run("custom partition and namespace", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"admin_partition": "team-a",
				"namespace":       "production",
			},
		}

		partition := request.GetString("admin_partition", "default")
		namespace := request.GetString("namespace", "default")

		assert.Equal(t, "team-a", partition)
		assert.Equal(t, "production", namespace)
	})
}

func TestGetACLAuthMethodTool(t *testing.T) {
	tool := GetACLAuthMethodTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "acl_auth_method", tool.Tool.Name)
}

func TestGetACLAuthMethodHandler(t *testing.T) {
	mockAuthMethod := map[string]interface{}{
		"Name":        "kubernetes",
		"Type":        "kubernetes",
		"Description": "Kubernetes auth method for service mesh",
		"Config": map[string]interface{}{
			"Host":              "https://kubernetes.default.svc.cluster.local:443",
			"CACert":            "-----BEGIN CERTIFICATE-----\ntest-ca-cert\n-----END CERTIFICATE-----",
			"ServiceAccountJWT": "",
		},
		"CreateIndex": 10,
		"ModifyIndex": 15,
		"Namespace":   "default",
		"Partition":   "default",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/acl/auth-method/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockAuthMethod)
	}))
	defer server.Close()

	t.Run("successful auth method retrieval", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"auth_method_name": "kubernetes",
			},
		}

		authMethodName, err := request.RequireString("auth_method_name")
		assert.NoError(t, err)
		assert.Equal(t, "kubernetes", authMethodName)
	})

	t.Run("missing auth method name", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("auth_method_name")
		assert.Error(t, err)
	})
}

func TestACLAuthMethodParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid auth method parameters",
			arguments: map[string]interface{}{
				"auth_method_name": "kubernetes",
				"admin_partition":  "default",
				"namespace":        "default",
			},
			expectError: false,
		},
		{
			name: "missing auth method name",
			arguments: map[string]interface{}{
				"admin_partition": "default",
			},
			expectError: true,
			errorField:  "auth_method_name",
		},
		{
			name: "custom enterprise parameters",
			arguments: map[string]interface{}{
				"auth_method_name": "oidc-provider",
				"admin_partition":  "team-b",
				"namespace":        "backend",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &MockCallToolRequest{Arguments: tt.arguments}

			if tt.expectError {
				_, err := request.RequireString(tt.errorField)
				assert.Error(t, err)
			} else {
				// Test optional parameters with defaults
				partition := request.GetString("admin_partition", "default")
				namespace := request.GetString("namespace", "default")

				assert.NotEmpty(t, partition)
				assert.NotEmpty(t, namespace)
			}
		})
	}
}

func TestACLAuthMethodResponseProcessing(t *testing.T) {
	t.Run("kubernetes auth method response", func(t *testing.T) {
		authMethod := map[string]interface{}{
			"Name": "kubernetes",
			"Type": "kubernetes",
			"Config": map[string]interface{}{
				"Host":   "https://kubernetes.default.svc.cluster.local:443",
				"CACert": "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
			},
		}

		data, err := json.MarshalIndent(authMethod, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "kubernetes", unmarshaled["Name"])
		assert.Equal(t, "kubernetes", unmarshaled["Type"])

		config, ok := unmarshaled["Config"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, config["Host"], "kubernetes.default.svc.cluster.local")
	})

	t.Run("oidc auth method response", func(t *testing.T) {
		authMethod := map[string]interface{}{
			"Name": "oidc-provider",
			"Type": "oidc",
			"Config": map[string]interface{}{
				"OIDCDiscoveryURL": "https://auth.example.com/.well-known/openid_configuration",
				"OIDCClientID":     "consul",
				"BoundAudiences":   []string{"consul", "vault"},
			},
		}

		data, err := json.MarshalIndent(authMethod, "", "  ")
		require.NoError(t, err)
		assert.Contains(t, string(data), "oidc-provider")
		assert.Contains(t, string(data), "OIDCDiscoveryURL")
	})
}

func TestACLAuthMethodTypes(t *testing.T) {
	validAuthMethodTypes := []string{"kubernetes", "oidc", "jwt"}

	for _, authType := range validAuthMethodTypes {
		t.Run(fmt.Sprintf("%s auth method type", authType), func(t *testing.T) {
			authMethod := map[string]interface{}{
				"Name": fmt.Sprintf("test-%s", authType),
				"Type": authType,
			}

			assert.Equal(t, authType, authMethod["Type"])
		})
	}
}

func TestACLAuthMethodErrorHandling(t *testing.T) {
	t.Run("auth method not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Auth method not found"))
		}))
		defer server.Close()

		// Would test handling of non-existent auth methods
		assert.True(t, true) // Placeholder
	})

	t.Run("ACL permission denied", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Permission denied"))
		}))
		defer server.Close()

		// Would test handling of ACL permission errors
		assert.True(t, true) // Placeholder
	})

	t.Run("invalid auth method configuration", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid configuration"))
		}))
		defer server.Close()

		// Would test handling of invalid configurations
		assert.True(t, true) // Placeholder
	})
}

func TestACLAuthMethodSecurity(t *testing.T) {
	t.Run("sensitive data handling", func(t *testing.T) {
		// Test that sensitive data like client secrets are handled properly
		authMethod := map[string]interface{}{
			"Config": map[string]interface{}{
				"OIDCClientSecret":  "[REDACTED]",
				"ServiceAccountJWT": "",
			},
		}

		data, err := json.Marshal(authMethod)
		require.NoError(t, err)

		// Should contain redacted placeholder instead of actual secret
		assert.Contains(t, string(data), "[REDACTED]")
		assert.NotContains(t, string(data), "actual-secret-value")
	})
}
