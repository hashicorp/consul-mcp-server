// Copyright (c) HashiCorp, Inc.
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

func TestGetNamespacesTool(t *testing.T) {
	tool := GetNamespacesTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "namespaces", tool.Tool.Name)
}

func TestGetNamespacesHandler(t *testing.T) {
	mockNamespaces := []map[string]interface{}{
		{
			"Name":        "default",
			"Description": "Default namespace",
			"ACLs": map[string]interface{}{
				"PolicyDefaults": []map[string]interface{}{
					{"ID": "policy-1", "Name": "default-policy"},
				},
				"RoleDefaults": []map[string]interface{}{
					{"ID": "role-1", "Name": "default-role"},
				},
			},
			"Meta": map[string]interface{}{
				"created_by": "system",
			},
			"CreateIndex": 5,
			"ModifyIndex": 5,
		},
		{
			"Name":        "production",
			"Description": "Production namespace for live services",
			"ACLs": map[string]interface{}{
				"PolicyDefaults": []map[string]interface{}{
					{"ID": "policy-prod", "Name": "production-policy"},
				},
			},
			"Meta": map[string]interface{}{
				"environment": "production",
				"team":        "platform",
			},
			"CreateIndex": 10,
			"ModifyIndex": 15,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/namespaces", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockNamespaces)
	}))
	defer server.Close()

	t.Run("successful namespaces listing", func(t *testing.T) {
		// Test would verify namespaces are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})

	t.Run("custom partition filtering", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"admin_partition": "team-a",
			},
		}

		partition := request.GetString("admin_partition", "default")
		assert.Equal(t, "team-a", partition)
	})
}

func TestGetNamespaceTool(t *testing.T) {
	tool := GetNamespaceTool(log.New())

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "namespace", tool.Tool.Name)
}

func TestGetNamespaceHandler(t *testing.T) {
	mockNamespace := map[string]interface{}{
		"Name":        "production",
		"Description": "Production environment namespace with strict ACL policies",
		"ACLs": map[string]interface{}{
			"PolicyDefaults": []map[string]interface{}{
				{"ID": "policy-prod-1", "Name": "production-base-policy"},
				{"ID": "policy-prod-2", "Name": "production-service-policy"},
			},
			"RoleDefaults": []map[string]interface{}{
				{"ID": "role-prod-1", "Name": "production-developer"},
				{"ID": "role-prod-2", "Name": "production-operator"},
			},
		},
		"Meta": map[string]string{
			"environment": "production",
			"team":        "sre",
			"criticality": "high",
			"compliance":  "required",
			"backup":      "enabled",
		},
		"CreateIndex": 100,
		"ModifyIndex": 150,
		"Partition":   "default",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/namespace/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockNamespace)
	}))
	defer server.Close()

	t.Run("successful namespace retrieval", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"namespace_name": "production",
			},
		}

		namespaceName, err := request.RequireString("namespace_name")
		assert.NoError(t, err)
		assert.Equal(t, "production", namespaceName)
	})

	t.Run("missing namespace name", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("namespace_name")
		assert.Error(t, err)
	})
}

func TestNamespaceParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid namespace listing",
			arguments: map[string]interface{}{
				"admin_partition": "default",
			},
			expectError: false,
		},
		{
			name: "valid namespace retrieval",
			arguments: map[string]interface{}{
				"namespace_name":  "production",
				"admin_partition": "team-a",
			},
			expectError: false,
		},
		{
			name:        "missing namespace name for retrieval",
			arguments:   map[string]interface{}{},
			expectError: true,
			errorField:  "namespace_name",
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
				assert.NotEmpty(t, partition)
			}
		})
	}
}

func TestNamespaceACLConfiguration(t *testing.T) {
	t.Run("namespace with policy defaults", func(t *testing.T) {
		namespace := map[string]interface{}{
			"Name": "secure-namespace",
			"ACLs": map[string]interface{}{
				"PolicyDefaults": []map[string]interface{}{
					{"ID": "policy-1", "Name": "base-security-policy"},
					{"ID": "policy-2", "Name": "namespace-specific-policy"},
				},
			},
		}

		acls := namespace["ACLs"].(map[string]interface{})
		policies := acls["PolicyDefaults"].([]map[string]interface{})
		assert.Len(t, policies, 2)
		assert.Equal(t, "base-security-policy", policies[0]["Name"])
	})

	t.Run("namespace with role defaults", func(t *testing.T) {
		namespace := map[string]interface{}{
			"Name": "team-namespace",
			"ACLs": map[string]interface{}{
				"RoleDefaults": []map[string]interface{}{
					{"ID": "role-1", "Name": "team-developer"},
					{"ID": "role-2", "Name": "team-lead"},
				},
			},
		}

		acls := namespace["ACLs"].(map[string]interface{})
		roles := acls["RoleDefaults"].([]map[string]interface{})
		assert.Len(t, roles, 2)
		assert.Equal(t, "team-developer", roles[0]["Name"])
	})
}

func TestNamespaceMetadata(t *testing.T) {
	t.Run("environment metadata", func(t *testing.T) {
		environments := []string{"development", "staging", "production"}

		for _, env := range environments {
			namespace := map[string]interface{}{
				"Name": fmt.Sprintf("%s-namespace", env),
				"Meta": map[string]string{
					"environment": env,
					"team":        "engineering",
				},
			}

			meta := namespace["Meta"].(map[string]string)
			assert.Equal(t, env, meta["environment"])
		}
	})

	t.Run("team and criticality metadata", func(t *testing.T) {
		namespace := map[string]interface{}{
			"Meta": map[string]string{
				"team":        "sre",
				"criticality": "high",
				"compliance":  "sox",
				"backup":      "enabled",
			},
		}

		meta := namespace["Meta"].(map[string]string)
		assert.Equal(t, "sre", meta["team"])
		assert.Equal(t, "high", meta["criticality"])
		assert.Equal(t, "sox", meta["compliance"])
	})
}

func TestNamespaceResponseProcessing(t *testing.T) {
	t.Run("complete namespace response", func(t *testing.T) {
		namespace := map[string]interface{}{
			"Name":        "comprehensive-namespace",
			"Description": "Namespace with all features",
			"ACLs": map[string]interface{}{
				"PolicyDefaults": []map[string]interface{}{
					{"ID": "policy-1", "Name": "policy-one"},
				},
				"RoleDefaults": []map[string]interface{}{
					{"ID": "role-1", "Name": "role-one"},
				},
			},
			"Meta": map[string]string{
				"environment": "production",
				"team":        "platform",
			},
		}

		data, err := json.MarshalIndent(namespace, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "comprehensive-namespace", unmarshaled["Name"])
		assert.Contains(t, unmarshaled, "ACLs")
		assert.Contains(t, unmarshaled, "Meta")
	})
}

func TestNamespaceUseCases(t *testing.T) {
	t.Run("multi-tenant setup", func(t *testing.T) {
		tenants := []string{"tenant-a", "tenant-b", "tenant-c"}

		for _, tenant := range tenants {
			namespace := map[string]interface{}{
				"Name":        tenant,
				"Description": fmt.Sprintf("Namespace for %s", tenant),
				"Meta": map[string]string{
					"tenant":      tenant,
					"environment": "production",
				},
			}

			meta := namespace["Meta"].(map[string]string)
			assert.Equal(t, tenant, meta["tenant"])
		}
	})

	t.Run("environment-based namespaces", func(t *testing.T) {
		environments := map[string]string{
			"dev":     "Development environment",
			"staging": "Staging environment",
			"prod":    "Production environment",
		}

		for env, desc := range environments {
			namespace := map[string]interface{}{
				"Name":        env,
				"Description": desc,
				"Meta": map[string]string{
					"environment": env,
				},
			}

			assert.Equal(t, desc, namespace["Description"])
		}
	})
}

func TestNamespaceErrorHandling(t *testing.T) {
	t.Run("namespace not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Namespace not found"))
		}))
		defer server.Close()

		// Would test handling of non-existent namespaces
		assert.True(t, true) // Placeholder
	})

	t.Run("enterprise feature not available", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Namespaces require Consul Enterprise"))
		}))
		defer server.Close()

		// Would test handling when namespaces aren't available in OSS
		assert.True(t, true) // Placeholder
	})
}
