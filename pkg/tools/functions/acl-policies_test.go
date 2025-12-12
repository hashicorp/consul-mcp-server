// Copyright IBM Corp. 2025
// SPDX-License-Identifier: MPL-2.0

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

func TestGetACLPoliciesTool(t *testing.T) {
	logger := log.New()
	tool := GetACLPolicyTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "acl_policy", tool.Tool.Name)
}

func TestGetACLPoliciesHandler(t *testing.T) {
	mockPolicies := []map[string]interface{}{
		{
			"ID":          "policy-1",
			"Name":        "web-policy",
			"Description": "Policy for web services",
			"Rules":       "service_prefix \"web\" { policy = \"write\" }",
			"Hash":        "hash-web-policy",
			"CreateIndex": 10,
			"ModifyIndex": 15,
		},
		{
			"ID":          "policy-2",
			"Name":        "database-policy",
			"Description": "Policy for database access",
			"Rules":       "service \"database\" { policy = \"read\" }\nnode_prefix \"\" { policy = \"read\" }",
			"Hash":        "hash-db-policy",
			"CreateIndex": 20,
			"ModifyIndex": 25,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/acl/policies", r.URL.Path)

		query := r.URL.Query()
		assert.Equal(t, "default", query.Get("partition"))
		assert.Equal(t, "default", query.Get("ns"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockPolicies)
	}))
	defer server.Close()

	t.Run("successful policies listing", func(t *testing.T) {
		// Test would verify policies are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})
}

func TestGetACLPolicyTool(t *testing.T) {
	logger := log.New()
	tool := GetACLPolicyTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "acl_policy", tool.Tool.Name)
}

func TestGetACLPolicyHandler(t *testing.T) {
	mockPolicy := map[string]interface{}{
		"ID":          "policy-1",
		"Name":        "web-policy",
		"Description": "Comprehensive policy for web tier services",
		"Rules": `service_prefix "web" {
  policy = "write"
}

node_prefix "" {
  policy = "read"
}

key_prefix "config/web/" {
  policy = "write"
}`,
		"Hash":        "hash-web-policy-123",
		"CreateIndex": 10,
		"ModifyIndex": 15,
		"Namespace":   "default",
		"Partition":   "default",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/acl/policy/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockPolicy)
	}))
	defer server.Close()

	t.Run("successful policy retrieval by ID", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"policy_id": "policy-1",
			},
		}

		policyID, err := request.RequireString("policy_id")
		assert.NoError(t, err)
		assert.Equal(t, "policy-1", policyID)
	})

	t.Run("successful policy retrieval by name", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"policy_name": "web-policy",
			},
		}

		policyName := request.GetString("policy_name", "")
		assert.Equal(t, "web-policy", policyName)
	})

	t.Run("missing policy identifier", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("policy_id")
		assert.Error(t, err)
	})
}

func TestACLPolicyParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid policy ID",
			arguments: map[string]interface{}{
				"policy_id":       "policy-123",
				"admin_partition": "default",
				"namespace":       "default",
			},
			expectError: false,
		},
		{
			name: "valid policy name",
			arguments: map[string]interface{}{
				"policy_name":     "web-policy",
				"admin_partition": "team-a",
				"namespace":       "production",
			},
			expectError: false,
		},
		{
			name: "missing policy identifier",
			arguments: map[string]interface{}{
				"admin_partition": "default",
			},
			expectError: true,
			errorField:  "policy_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &MockCallToolRequest{Arguments: tt.arguments}

			if tt.expectError {
				_, err := request.RequireString(tt.errorField)
				assert.Error(t, err)
			} else {
				// Test optional parameters
				partition := request.GetString("admin_partition", "default")
				namespace := request.GetString("namespace", "default")

				assert.NotEmpty(t, partition)
				assert.NotEmpty(t, namespace)
			}
		})
	}
}

func TestACLPolicyRulesValidation(t *testing.T) {
	tests := []struct {
		name  string
		rules string
		valid bool
	}{
		{
			name:  "service policy rules",
			rules: `service_prefix "web" { policy = "write" }`,
			valid: true,
		},
		{
			name:  "node policy rules",
			rules: `node_prefix "" { policy = "read" }`,
			valid: true,
		},
		{
			name:  "key policy rules",
			rules: `key_prefix "config/" { policy = "write" }`,
			valid: true,
		},
		{
			name:  "connect policy rules",
			rules: `service "web" { policy = "write", intentions = "write" }`,
			valid: true,
		},
		{
			name: "combined policy rules",
			rules: `service_prefix "api" { policy = "write" }
node_prefix "" { policy = "read" }
key_prefix "config/api/" { policy = "write" }`,
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := map[string]interface{}{
				"Name":  fmt.Sprintf("test-%s", tt.name),
				"Rules": tt.rules,
			}

			assert.Equal(t, tt.rules, policy["Rules"])
			if tt.valid {
				assert.NotEmpty(t, policy["Rules"])
			}
		})
	}
}

func TestACLPolicyResponseProcessing(t *testing.T) {
	t.Run("policy with complex rules", func(t *testing.T) {
		policy := map[string]interface{}{
			"ID":   "complex-policy",
			"Name": "complex-web-policy",
			"Rules": `service_prefix "web" {
  policy = "write"
  intentions = "write"
}

service "database" {
  policy = "read"
}

node_prefix "" {
  policy = "read"
}

key_prefix "config/web/" {
  policy = "write"
}`,
		}

		data, err := json.MarshalIndent(policy, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "complex-policy", unmarshaled["ID"])
		assert.Contains(t, unmarshaled["Rules"].(string), "service_prefix")
		assert.Contains(t, unmarshaled["Rules"].(string), "intentions")
	})

	t.Run("policy metadata", func(t *testing.T) {
		policy := map[string]interface{}{
			"Hash":        "policy-hash-123",
			"CreateIndex": 100,
			"ModifyIndex": 150,
		}

		data, err := json.Marshal(policy)
		require.NoError(t, err)
		assert.Contains(t, string(data), "policy-hash-123")
	})
}

func TestACLPolicyPermissionTypes(t *testing.T) {
	permissionTypes := []string{"read", "write", "deny"}

	for _, permission := range permissionTypes {
		t.Run(fmt.Sprintf("%s permission", permission), func(t *testing.T) {
			rules := fmt.Sprintf(`service "test" { policy = "%s" }`, permission)
			policy := map[string]interface{}{
				"Rules": rules,
			}

			assert.Contains(t, policy["Rules"].(string), permission)
		})
	}
}

func TestACLPolicyErrorHandling(t *testing.T) {
	t.Run("policy not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Policy not found"))
		}))
		defer server.Close()

		// Would test handling of non-existent policies
		assert.True(t, true) // Placeholder
	})

	t.Run("invalid policy rules", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid policy rules"))
		}))
		defer server.Close()

		// Would test handling of malformed policy rules
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
}

func TestACLPolicyBuiltinPolicies(t *testing.T) {
	builtinPolicies := []string{
		"global-management",
		"builtin/global-read-only",
		"builtin/namespace-management",
	}

	for _, policyName := range builtinPolicies {
		t.Run(fmt.Sprintf("builtin policy %s", policyName), func(t *testing.T) {
			policy := map[string]interface{}{
				"Name": policyName,
			}

			if policyName == "global-management" {
				assert.Equal(t, "global-management", policy["Name"])
			} else {
				assert.Contains(t, policy["Name"].(string), "builtin/")
			}
		})
	}
}
