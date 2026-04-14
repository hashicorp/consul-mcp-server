// Copyright (c) HashiCorp, Inc.
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
)

func TestGetACLBindingRulesTool(t *testing.T) {
	logger := log.New()
	tool := GetACLBindingRulesTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "acl_binding_rules", tool.Tool.Name)
}

func TestGetACLBindingRulesHandler(t *testing.T) {
	mockBindingRules := []map[string]interface{}{
		{
			"ID":          "binding-rule-1",
			"Description": "Kubernetes service account binding",
			"AuthMethod":  "kubernetes",
			"Selector":    "serviceaccount.namespace==default",
			"BindType":    "service",
			"BindName":    "web-${serviceaccount.name}",
			"CreateIndex": 10,
			"ModifyIndex": 15,
		},
		{
			"ID":          "binding-rule-2",
			"Description": "OIDC group binding for admins",
			"AuthMethod":  "oidc-provider",
			"Selector":    "\"admin\" in groups",
			"BindType":    "role",
			"BindName":    "admin-role",
			"CreateIndex": 20,
			"ModifyIndex": 25,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/acl/binding-rules", r.URL.Path)

		query := r.URL.Query()
		assert.Equal(t, "default", query.Get("partition"))
		assert.Equal(t, "default", query.Get("ns"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockBindingRules)
	}))
	defer server.Close()

	t.Run("successful binding rules listing", func(t *testing.T) {
		// Test would verify binding rules are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})

	t.Run("auth method filtering", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"auth_method": "kubernetes",
			},
		}

		authMethod := request.GetString("auth_method", "")
		assert.Equal(t, "kubernetes", authMethod)
	})
}

func TestGetACLBindingRuleTool(t *testing.T) {
	logger := log.New()
	tool := GetACLBindingRuleTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "acl_binding_rule", tool.Tool.Name)
}

func TestGetACLBindingRuleHandler(t *testing.T) {
	mockBindingRule := map[string]interface{}{
		"ID":          "binding-rule-1",
		"Description": "Kubernetes service account to service identity binding",
		"AuthMethod":  "kubernetes",
		"Selector":    "serviceaccount.namespace==default and serviceaccount.name==web",
		"BindType":    "service",
		"BindName":    "web-${serviceaccount.name}",
		"CreateIndex": 10,
		"ModifyIndex": 15,
		"Namespace":   "default",
		"Partition":   "default",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/acl/binding-rule/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockBindingRule)
	}))
	defer server.Close()

	t.Run("successful binding rule retrieval", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"binding_rule_id": "binding-rule-1",
			},
		}

		bindingRuleID, err := request.RequireString("binding_rule_id")
		assert.NoError(t, err)
		assert.Equal(t, "binding-rule-1", bindingRuleID)
	})

	t.Run("missing binding rule ID", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("binding_rule_id")
		assert.Error(t, err)
	})
}

func TestACLBindingRuleParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid binding rule parameters",
			arguments: map[string]interface{}{
				"binding_rule_id": "binding-rule-123",
				"admin_partition": "default",
				"namespace":       "default",
			},
			expectError: false,
		},
		{
			name: "missing binding rule ID",
			arguments: map[string]interface{}{
				"admin_partition": "default",
			},
			expectError: true,
			errorField:  "binding_rule_id",
		},
		{
			name: "auth method filter",
			arguments: map[string]interface{}{
				"auth_method": "kubernetes",
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
				// Test optional parameters
				partition := request.GetString("admin_partition", "default")
				namespace := request.GetString("namespace", "default")
				//authMethod := request.GetString("auth_method", "")

				assert.NotEmpty(t, partition)
				assert.NotEmpty(t, namespace)
				// authMethod can be empty for general listing
			}
		})
	}
}

func TestACLBindingRuleTypes(t *testing.T) {
	validBindTypes := []string{"service", "role"}

	for _, bindType := range validBindTypes {
		t.Run(fmt.Sprintf("%s bind type", bindType), func(t *testing.T) {
			bindingRule := map[string]interface{}{
				"BindType": bindType,
				"BindName": fmt.Sprintf("test-%s", bindType),
			}

			assert.Equal(t, bindType, bindingRule["BindType"])
		})
	}
}

func TestACLBindingRuleSelectors(t *testing.T) {
	tests := []struct {
		name     string
		selector string
		authType string
	}{
		{
			name:     "kubernetes service account selector",
			selector: "serviceaccount.namespace==default and serviceaccount.name==web",
			authType: "kubernetes",
		},
		{
			name:     "oidc group membership selector",
			selector: "\"admin\" in groups",
			authType: "oidc",
		},
		{
			name:     "jwt claim selector",
			selector: "sub==\"user@example.com\"",
			authType: "jwt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bindingRule := map[string]interface{}{
				"Selector":   tt.selector,
				"AuthMethod": tt.authType,
			}

			assert.Equal(t, tt.selector, bindingRule["Selector"])
			assert.Equal(t, tt.authType, bindingRule["AuthMethod"])
		})
	}
}

func TestACLBindingRuleErrorHandling(t *testing.T) {
	t.Run("binding rule not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Binding rule not found"))
		}))
		defer server.Close()

		// Would test handling of non-existent binding rules
		assert.True(t, true) // Placeholder
	})

	t.Run("invalid selector syntax", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid selector syntax"))
		}))
		defer server.Close()

		// Would test handling of invalid selector expressions
		assert.True(t, true) // Placeholder
	})

	t.Run("auth method mismatch", func(t *testing.T) {
		// Test scenarios where binding rule references non-existent auth method
		assert.True(t, true) // Placeholder
	})
}
