// filepath: /Users/srahul3/git/consul-mcp-server/pkg/tools/functions/test_helpers.go
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functions

import "fmt"

// MockCallToolRequest is a shared mock implementation for testing
type MockCallToolRequest struct {
	Arguments map[string]interface{}
}

func (m *MockCallToolRequest) GetString(key, defaultValue string) string {
	if val, ok := m.Arguments[key].(string); ok {
		return val
	}
	return defaultValue
}

func (m *MockCallToolRequest) RequireString(key string) (string, error) {
	if val, ok := m.Arguments[key].(string); ok && val != "" {
		return val, nil
	}
	return "", fmt.Errorf("required argument %q not found", key)
}

func (m *MockCallToolRequest) GetArguments() map[string]interface{} {
	return m.Arguments
}

// Services represents a map of service names to their tags
type Services map[string][]string
