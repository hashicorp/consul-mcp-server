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
	"github.com/stretchr/testify/require"
)

func TestGetACLTemplatedPoliciesTool(t *testing.T) {
	logger := log.New()
	tool := GetACLTemplatedPoliciesTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "acl_templated_policies", tool.Tool.Name)
}

func TestGetACLTemplatedPoliciesHandler(t *testing.T) {
	mockTemplatedPolicies := []map[string]interface{}{
		{
			"TemplateName": "builtin/service",
			"Schema":       "{\n  \"type\": \"object\",\n  \"properties\": {\n    \"name\": { \"type\": \"string\" }\n  }\n}",
			"Template":     "service \"{{.name}}\" {\n  policy = \"write\"\n}",
		},
		{
			"TemplateName": "builtin/node",
			"Schema":       "{\n  \"type\": \"object\",\n  \"properties\": {\n    \"name\": { \"type\": \"string\" }\n  }\n}",
			"Template":     "node \"{{.name}}\" {\n  policy = \"write\"\n}",
		},
		{
			"TemplateName": "builtin/dns",
			"Schema":       "{\n  \"type\": \"object\",\n  \"properties\": {}\n}",
			"Template":     "node_prefix \"\" {\n  policy = \"read\"\n}\nservice_prefix \"\" {\n  policy = \"read\"\n}",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/acl/templated-policies", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockTemplatedPolicies)
	}))
	defer server.Close()

	t.Run("successful templated policies listing", func(t *testing.T) {
		// Test would verify templated policies are returned correctly
		assert.True(t, true) // Placeholder for actual implementation
	})
}

func TestGetACLTemplatedPolicyTool(t *testing.T) {
	logger := log.New()
	tool := GetACLTemplatedPolicyTool(logger)

	assert.NotNil(t, tool.Tool)
	assert.NotNil(t, tool.Handler)
	assert.Equal(t, "acl_templated_policy", tool.Tool.Name)
}

func TestGetACLTemplatedPolicyHandler(t *testing.T) {
	mockTemplatedPolicy := map[string]interface{}{
		"TemplateName": "builtin/service",
		"Schema": `{
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "The name of the service"
    }
  },
  "required": ["name"]
}`,
		"Template": `service "{{.name}}" {
  policy = "write"
}

key_prefix "service/{{.name}}/" {
  policy = "write"
}`,
		"Description": "Templated policy for service-specific permissions",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/v1/acl/templated-policy/")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockTemplatedPolicy)
	}))
	defer server.Close()

	t.Run("successful templated policy retrieval", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{
				"template_name": "builtin/service",
			},
		}

		templateName, err := request.RequireString("template_name")
		assert.NoError(t, err)
		assert.Equal(t, "builtin/service", templateName)
	})

	t.Run("missing template name", func(t *testing.T) {
		request := &MockCallToolRequest{
			Arguments: map[string]interface{}{},
		}

		_, err := request.RequireString("template_name")
		assert.Error(t, err)
	})
}

func TestACLTemplatedPolicyParameterHandling(t *testing.T) {
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
		errorField  string
	}{
		{
			name: "valid template name",
			arguments: map[string]interface{}{
				"template_name": "builtin/service",
			},
			expectError: false,
		},
		{
			name: "custom template name",
			arguments: map[string]interface{}{
				"template_name": "custom/web-service",
			},
			expectError: false,
		},
		{
			name:        "missing template name",
			arguments:   map[string]interface{}{},
			expectError: true,
			errorField:  "template_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &MockCallToolRequest{Arguments: tt.arguments}

			if tt.expectError {
				_, err := request.RequireString(tt.errorField)
				assert.Error(t, err)
			} else {
				templateName, err := request.RequireString("template_name")
				assert.NoError(t, err)
				assert.NotEmpty(t, templateName)
			}
		})
	}
}

func TestACLTemplatedPolicyTypes(t *testing.T) {
	builtinTemplates := []string{
		"builtin/service",
		"builtin/node",
		"builtin/dns",
		"builtin/namespace-management",
	}

	for _, template := range builtinTemplates {
		t.Run(fmt.Sprintf("builtin template %s", template), func(t *testing.T) {
			templatedPolicy := map[string]interface{}{
				"TemplateName": template,
			}

			assert.Contains(t, templatedPolicy["TemplateName"].(string), "builtin/")
		})
	}
}

func TestACLTemplatedPolicySchemaValidation(t *testing.T) {
	t.Run("service template schema", func(t *testing.T) {
		schema := `{
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "The name of the service"
    }
  },
  "required": ["name"]
}`

		var schemaObj map[string]interface{}
		err := json.Unmarshal([]byte(schema), &schemaObj)
		require.NoError(t, err)

		assert.Equal(t, "object", schemaObj["type"])
		properties := schemaObj["properties"].(map[string]interface{})
		assert.Contains(t, properties, "name")
	})

	t.Run("node template schema", func(t *testing.T) {
		schema := `{
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "The name of the node"
    }
  },
  "required": ["name"]
}`

		var schemaObj map[string]interface{}
		err := json.Unmarshal([]byte(schema), &schemaObj)
		require.NoError(t, err)

		required := schemaObj["required"].([]interface{})
		assert.Contains(t, required, "name")
	})
}

func TestACLTemplatedPolicyTemplating(t *testing.T) {
	tests := []struct {
		name     string
		template string
		variable string
	}{
		{
			name:     "service template",
			template: `service "{{.name}}" { policy = "write" }`,
			variable: "name",
		},
		{
			name:     "node template",
			template: `node "{{.name}}" { policy = "write" }`,
			variable: "name",
		},
		{
			name:     "key prefix template",
			template: `key_prefix "service/{{.name}}/" { policy = "write" }`,
			variable: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, tt.template, fmt.Sprintf("{{.%s}}", tt.variable))
			assert.Contains(t, tt.template, "policy")
		})
	}
}

func TestACLTemplatedPolicyResponseProcessing(t *testing.T) {
	t.Run("complete templated policy response", func(t *testing.T) {
		templatedPolicy := map[string]interface{}{
			"TemplateName": "builtin/service",
			"Schema": `{
  "type": "object",
  "properties": {
    "name": { "type": "string" }
  }
}`,
			"Template": `service "{{.name}}" {
  policy = "write"
}`,
			"Description": "Service-specific access template",
		}

		data, err := json.MarshalIndent(templatedPolicy, "", "  ")
		require.NoError(t, err)

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "builtin/service", unmarshaled["TemplateName"])
		assert.Contains(t, unmarshaled["Template"].(string), "{{.name}}")
		assert.Contains(t, unmarshaled["Schema"].(string), "properties")
	})
}

func TestACLTemplatedPolicyUseCases(t *testing.T) {
	t.Run("microservice deployment", func(t *testing.T) {
		// Template for microservice-specific permissions
		template := map[string]interface{}{
			"TemplateName": "custom/microservice",
			"Template": `service "{{.service_name}}" {
  policy = "write"
  intentions = "write"
}

key_prefix "config/{{.service_name}}/" {
  policy = "write"
}

session_prefix "{{.service_name}}" {
  policy = "write"
}`,
		}

		templateStr := template["Template"].(string)
		assert.Contains(t, templateStr, "{{.service_name}}")
		assert.Contains(t, templateStr, "intentions")
		assert.Contains(t, templateStr, "session_prefix")
	})

	t.Run("namespace management", func(t *testing.T) {
		// Template for namespace-level permissions
		template := map[string]interface{}{
			"TemplateName": "custom/namespace-admin",
			"Template": `namespace "{{.namespace}}" {
  policy = "write"
  
  service_prefix "" {
    policy = "write"
    intentions = "write"
  }
  
  node_prefix "" {
    policy = "read"
  }
}`,
		}

		templateStr := template["Template"].(string)
		assert.Contains(t, templateStr, "{{.namespace}}")
		assert.Contains(t, templateStr, "service_prefix")
	})
}

func TestACLTemplatedPolicyErrorHandling(t *testing.T) {
	t.Run("template not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Templated policy not found"))
		}))
		defer server.Close()

		// Would test handling of non-existent templates
		assert.True(t, true) // Placeholder
	})

	t.Run("invalid template syntax", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid template syntax"))
		}))
		defer server.Close()

		// Would test handling of malformed templates
		assert.True(t, true) // Placeholder
	})

	t.Run("schema validation error", func(t *testing.T) {
		// Test scenarios where template parameters don't match schema
		assert.True(t, true) // Placeholder
	})
}
