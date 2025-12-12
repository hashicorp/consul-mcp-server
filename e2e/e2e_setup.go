// Copyright IBM Corp. 2025
// SPDX-License-Identifier: MPL-2.0

package e2e

type ContentType string

const (
	CONST_TYPE_RESOURCE    ContentType = "resources"
	CONST_TYPE_DATA_SOURCE ContentType = "data-sources"
	CONST_TYPE_GUIDES      ContentType = "guides"
	CONST_TYPE_FUNCTIONS   ContentType = "functions"
	CONST_TYPE_OVERVIEW    ContentType = "overview"
)

type TestCase struct {
	TestName        string                 `json:"testName"`
	TestShouldFail  bool                   `json:"testShouldFail"`
	TestDescription string                 `json:"testDescription"`
	TestContentType ContentType            `json:"testContentType,omitempty"`
	TestPayload     map[string]interface{} `json:"testPayload,omitempty"`
}

var toolsTestCases = []TestCase{
	{
		TestName:        "overview_documentation",
		TestShouldFail:  false,
		TestDescription: "Testing search_providers overview documentation with v2 API",
		TestContentType: CONST_TYPE_OVERVIEW,
		TestPayload: map[string]interface{}{
			"provider_name":      "google",
			"provider_namespace": "hashicorp",
			"provider_version":   "latest",
			"provider_data_type": "overview",
			"service_slug":       "index",
		},
	},
}

var resourcesTestCases = []TestCase{
	{
		TestName:        "empty_doc_id",
		TestShouldFail:  true,
		TestDescription: "Testing get_provider_details with empty provider_doc_id",
		TestPayload: map[string]interface{}{
			"provider_doc_id": "",
		},
	},
}
