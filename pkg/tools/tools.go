// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package tools

import (
	tools "github.com/hashicorp/consul-mcp-server/pkg/tools/functions"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
)

// MCPServerInterface defines the interface for MCP server operations needed by tools
type MCPServerInterface interface {
	AddTool(tool mcp.Tool, handler server.ToolHandlerFunc)
}

func RegisterTools(hcServer MCPServerInterface, logger *log.Logger) {
	// Consul Catalog tools
	getCatalogServicesTool := tools.GetCatalogServicesTool(logger)
	hcServer.AddTool(getCatalogServicesTool.Tool, getCatalogServicesTool.Handler)

	getCatalogNodesTool := tools.GetCatalogNodesTool(logger)
	hcServer.AddTool(getCatalogNodesTool.Tool, getCatalogNodesTool.Handler)

	getCatalogServiceTool := tools.GetCatalogServiceTool(logger)
	hcServer.AddTool(getCatalogServiceTool.Tool, getCatalogServiceTool.Handler)

	getCatalogConnectTool := tools.GetCatalogConnectTool(logger)
	hcServer.AddTool(getCatalogConnectTool.Tool, getCatalogConnectTool.Handler)

	getCatalogNodeTool := tools.GetCatalogNodeTool(logger)
	hcServer.AddTool(getCatalogNodeTool.Tool, getCatalogNodeTool.Handler)

	getCatalogDatacentersTool := tools.GetCatalogDatacentersTool(logger)
	hcServer.AddTool(getCatalogDatacentersTool.Tool, getCatalogDatacentersTool.Handler)

	getCatalogGatewayServicesTool := tools.GetCatalogGatewayServicesTool(logger)
	hcServer.AddTool(getCatalogGatewayServicesTool.Tool, getCatalogGatewayServicesTool.Handler)

	// Consul Agent Overview tools
	getAgentSelfTool := tools.GetAgentSelfTool(logger)
	hcServer.AddTool(getAgentSelfTool.Tool, getAgentSelfTool.Handler)

	getAgentConfigTool := tools.GetAgentConfigTool(logger)
	hcServer.AddTool(getAgentConfigTool.Tool, getAgentConfigTool.Handler)

	getAgentMembersTool := tools.GetAgentMembersTool(logger)
	hcServer.AddTool(getAgentMembersTool.Tool, getAgentMembersTool.Handler)

	getAgentMetricsTool := tools.GetAgentMetricsTool(logger)
	hcServer.AddTool(getAgentMetricsTool.Tool, getAgentMetricsTool.Handler)

	getAgentHostTool := tools.GetAgentHostTool(logger)
	hcServer.AddTool(getAgentHostTool.Tool, getAgentHostTool.Handler)

	getAgentVersionTool := tools.GetAgentVersionTool(logger)
	hcServer.AddTool(getAgentVersionTool.Tool, getAgentVersionTool.Handler)

	getAgentReloadTool := tools.GetAgentReloadTool(logger)
	hcServer.AddTool(getAgentReloadTool.Tool, getAgentReloadTool.Handler)

	// Consul Agent Check tools
	getAgentChecksTool := tools.GetAgentChecksTool(logger)
	hcServer.AddTool(getAgentChecksTool.Tool, getAgentChecksTool.Handler)

	getAgentCheckTool := tools.GetAgentCheckTool(logger)
	hcServer.AddTool(getAgentCheckTool.Tool, getAgentCheckTool.Handler)

	getAgentCheckByNameTool := tools.GetAgentCheckByNameTool(logger)
	hcServer.AddTool(getAgentCheckByNameTool.Tool, getAgentCheckByNameTool.Handler)

	getAgentServiceHealthTool := tools.GetAgentServiceHealthTool(logger)
	hcServer.AddTool(getAgentServiceHealthTool.Tool, getAgentServiceHealthTool.Handler)

	getAgentServiceHealthByNameTool := tools.GetAgentServiceHealthByNameTool(logger)
	hcServer.AddTool(getAgentServiceHealthByNameTool.Tool, getAgentServiceHealthByNameTool.Handler)

	// Consul Agent Service tools
	getAgentServicesTool := tools.GetAgentServicesTool(logger)
	hcServer.AddTool(getAgentServicesTool.Tool, getAgentServicesTool.Handler)

	getAgentServiceTool := tools.GetAgentServiceTool(logger)
	hcServer.AddTool(getAgentServiceTool.Tool, getAgentServiceTool.Handler)

	getAgentServiceConfigurationTool := tools.GetAgentServiceConfigurationTool(logger)
	hcServer.AddTool(getAgentServiceConfigurationTool.Tool, getAgentServiceConfigurationTool.Handler)

	getAgentServiceLocalStateTool := tools.GetAgentServiceLocalStateTool(logger)
	hcServer.AddTool(getAgentServiceLocalStateTool.Tool, getAgentServiceLocalStateTool.Handler)

	// Consul Agent Connect tools
	getAgentConnectCATool := tools.GetAgentConnectCATool(logger)
	hcServer.AddTool(getAgentConnectCATool.Tool, getAgentConnectCATool.Handler)

	getAgentConnectAuthorizeTool := tools.GetAgentConnectAuthorizeTool(logger)
	hcServer.AddTool(getAgentConnectAuthorizeTool.Tool, getAgentConnectAuthorizeTool.Handler)

	getAgentConnectProxyConfigTool := tools.GetAgentConnectProxyConfigTool(logger)
	hcServer.AddTool(getAgentConnectProxyConfigTool.Tool, getAgentConnectProxyConfigTool.Handler)

	getAgentConnectLeafCertTool := tools.GetAgentConnectLeafCertTool(logger)
	hcServer.AddTool(getAgentConnectLeafCertTool.Tool, getAgentConnectLeafCertTool.Handler)

	// Consul Peering tools
	getPeeringsTool := tools.GetPeeringsTool(logger)
	hcServer.AddTool(getPeeringsTool.Tool, getPeeringsTool.Handler)

	getPeeringTool := tools.GetPeeringTool(logger)
	hcServer.AddTool(getPeeringTool.Tool, getPeeringTool.Handler)

	getPeeringExportedServicesTool := tools.GetPeeringExportedServicesTool(logger)
	hcServer.AddTool(getPeeringExportedServicesTool.Tool, getPeeringExportedServicesTool.Handler)

	// Consul Config tools
	getConfigEntriesTool := tools.GetConfigEntriesTool(logger)
	hcServer.AddTool(getConfigEntriesTool.Tool, getConfigEntriesTool.Handler)

	getConfigEntryTool := tools.GetConfigEntryTool(logger)
	hcServer.AddTool(getConfigEntryTool.Tool, getConfigEntryTool.Handler)

	// Consul Connect CA tools
	getConnectCARootsTool := tools.GetConnectCARootsTool(logger)
	hcServer.AddTool(getConnectCARootsTool.Tool, getConnectCARootsTool.Handler)

	getConnectCAConfigurationTool := tools.GetConnectCAConfigurationTool(logger)
	hcServer.AddTool(getConnectCAConfigurationTool.Tool, getConnectCAConfigurationTool.Handler)

	// Consul Connect Intentions tools
	getConnectIntentionsTool := tools.GetConnectIntentionsTool(logger)
	hcServer.AddTool(getConnectIntentionsTool.Tool, getConnectIntentionsTool.Handler)

	getConnectIntentionTool := tools.GetConnectIntentionTool(logger)
	hcServer.AddTool(getConnectIntentionTool.Tool, getConnectIntentionTool.Handler)

	getConnectIntentionMatchTool := tools.GetConnectIntentionMatchTool(logger)
	hcServer.AddTool(getConnectIntentionMatchTool.Tool, getConnectIntentionMatchTool.Handler)

	getConnectIntentionCheckTool := tools.GetConnectIntentionCheckTool(logger)
	hcServer.AddTool(getConnectIntentionCheckTool.Tool, getConnectIntentionCheckTool.Handler)

	// Consul Discovery Chain tools
	getDiscoveryChainTool := tools.GetDiscoveryChainTool(logger)
	hcServer.AddTool(getDiscoveryChainTool.Tool, getDiscoveryChainTool.Handler)

	// Consul Exported Services tools
	getExportedServicesTool := tools.GetExportedServicesTool(logger)
	hcServer.AddTool(getExportedServicesTool.Tool, getExportedServicesTool.Handler)

	// Consul Health tools
	getHealthNodeTool := tools.GetHealthNodeTool(logger)
	hcServer.AddTool(getHealthNodeTool.Tool, getHealthNodeTool.Handler)

	getHealthChecksTool := tools.GetHealthChecksTool(logger)
	hcServer.AddTool(getHealthChecksTool.Tool, getHealthChecksTool.Handler)

	getHealthServiceTool := tools.GetHealthServiceTool(logger)
	hcServer.AddTool(getHealthServiceTool.Tool, getHealthServiceTool.Handler)

	getHealthConnectTool := tools.GetHealthConnectTool(logger)
	hcServer.AddTool(getHealthConnectTool.Tool, getHealthConnectTool.Handler)

	getHealthIngressTool := tools.GetHealthIngressTool(logger)
	hcServer.AddTool(getHealthIngressTool.Tool, getHealthIngressTool.Handler)

	getHealthStateTool := tools.GetHealthStateTool(logger)
	hcServer.AddTool(getHealthStateTool.Tool, getHealthStateTool.Handler)

	// Consul KV tools
	getKVTool := tools.GetKVTool(logger)
	hcServer.AddTool(getKVTool.Tool, getKVTool.Handler)

	getKVKeysTool := tools.GetKVKeysTool(logger)
	hcServer.AddTool(getKVKeysTool.Tool, getKVKeysTool.Handler)

	getKVRecursiveTool := tools.GetKVRecursiveTool(logger)
	hcServer.AddTool(getKVRecursiveTool.Tool, getKVRecursiveTool.Handler)

	// Consul Operator tools
	getOperatorAreasTool := tools.GetOperatorAreasTool(logger)
	hcServer.AddTool(getOperatorAreasTool.Tool, getOperatorAreasTool.Handler)

	getOperatorAreaTool := tools.GetOperatorAreaTool(logger)
	hcServer.AddTool(getOperatorAreaTool.Tool, getOperatorAreaTool.Handler)

	getOperatorAreaMembersTool := tools.GetOperatorAreaMembersTool(logger)
	hcServer.AddTool(getOperatorAreaMembersTool.Tool, getOperatorAreaMembersTool.Handler)

	getOperatorAutopilotConfigurationTool := tools.GetOperatorAutopilotConfigurationTool(logger)
	hcServer.AddTool(getOperatorAutopilotConfigurationTool.Tool, getOperatorAutopilotConfigurationTool.Handler)

	getOperatorAutopilotHealthTool := tools.GetOperatorAutopilotHealthTool(logger)
	hcServer.AddTool(getOperatorAutopilotHealthTool.Tool, getOperatorAutopilotHealthTool.Handler)

	getOperatorAutopilotStateTool := tools.GetOperatorAutopilotStateTool(logger)
	hcServer.AddTool(getOperatorAutopilotStateTool.Tool, getOperatorAutopilotStateTool.Handler)

	getOperatorKeyringTool := tools.GetOperatorKeyringTool(logger)
	hcServer.AddTool(getOperatorKeyringTool.Tool, getOperatorKeyringTool.Handler)

	getOperatorLicenseTool := tools.GetOperatorLicenseTool(logger)
	hcServer.AddTool(getOperatorLicenseTool.Tool, getOperatorLicenseTool.Handler)

	getOperatorRaftConfigurationTool := tools.GetOperatorRaftConfigurationTool(logger)
	hcServer.AddTool(getOperatorRaftConfigurationTool.Tool, getOperatorRaftConfigurationTool.Handler)

	getOperatorSegmentTool := tools.GetOperatorSegmentTool(logger)
	hcServer.AddTool(getOperatorSegmentTool.Tool, getOperatorSegmentTool.Handler)

	getOperatorUsageTool := tools.GetOperatorUsageTool(logger)
	hcServer.AddTool(getOperatorUsageTool.Tool, getOperatorUsageTool.Handler)

	// Consul Namespaces tools
	getNamespacesTool := tools.GetNamespacesTool(logger)
	hcServer.AddTool(getNamespacesTool.Tool, getNamespacesTool.Handler)

	getNamespaceTool := tools.GetNamespaceTool(logger)
	hcServer.AddTool(getNamespaceTool.Tool, getNamespaceTool.Handler)

	// Consul Query tools
	getQueryTool := tools.GetQueryTool(logger)
	hcServer.AddTool(getQueryTool.Tool, getQueryTool.Handler)

	getQueryByIdTool := tools.GetQueryByIdTool(logger)
	hcServer.AddTool(getQueryByIdTool.Tool, getQueryByIdTool.Handler)

	getQueryExecuteTool := tools.GetQueryExecuteTool(logger)
	hcServer.AddTool(getQueryExecuteTool.Tool, getQueryExecuteTool.Handler)

	getQueryExplainTool := tools.GetQueryExplainTool(logger)
	hcServer.AddTool(getQueryExplainTool.Tool, getQueryExplainTool.Handler)

	// Consul Session tools
	getSessionTool := tools.GetSessionTool(logger)
	hcServer.AddTool(getSessionTool.Tool, getSessionTool.Handler)

	getSessionNodeTool := tools.GetSessionNodeTool(logger)
	hcServer.AddTool(getSessionNodeTool.Tool, getSessionNodeTool.Handler)

	getSessionListTool := tools.GetSessionListTool(logger)
	hcServer.AddTool(getSessionListTool.Tool, getSessionListTool.Handler)

	// Consul Status tools
	getStatusLeaderTool := tools.GetStatusLeaderTool(logger)
	hcServer.AddTool(getStatusLeaderTool.Tool, getStatusLeaderTool.Handler)

	getStatusPeersTool := tools.GetStatusPeersTool(logger)
	hcServer.AddTool(getStatusPeersTool.Tool, getStatusPeersTool.Handler)

	// Consul ACL tools
	getACLTokensTool := tools.GetACLTokensTool(logger)
	hcServer.AddTool(getACLTokensTool.Tool, getACLTokensTool.Handler)

	getACLTokenTool := tools.GetACLTokenTool(logger)
	hcServer.AddTool(getACLTokenTool.Tool, getACLTokenTool.Handler)

	getACLTokenSelfTool := tools.GetACLTokenSelfTool(logger)
	hcServer.AddTool(getACLTokenSelfTool.Tool, getACLTokenSelfTool.Handler)

	getACLPolicesTool := tools.GetACLPolicesTool(logger)
	hcServer.AddTool(getACLPolicesTool.Tool, getACLPolicesTool.Handler)

	getACLPolicyTool := tools.GetACLPolicyTool(logger)
	hcServer.AddTool(getACLPolicyTool.Tool, getACLPolicyTool.Handler)

	getACLTemplatedPoliciesTool := tools.GetACLTemplatedPoliciesTool(logger)
	hcServer.AddTool(getACLTemplatedPoliciesTool.Tool, getACLTemplatedPoliciesTool.Handler)

	getACLTemplatedPolicyTool := tools.GetACLTemplatedPolicyTool(logger)
	hcServer.AddTool(getACLTemplatedPolicyTool.Tool, getACLTemplatedPolicyTool.Handler)

	getACLRolesTool := tools.GetACLRolesTool(logger)
	hcServer.AddTool(getACLRolesTool.Tool, getACLRolesTool.Handler)

	getACLRoleTool := tools.GetACLRoleTool(logger)
	hcServer.AddTool(getACLRoleTool.Tool, getACLRoleTool.Handler)

	getACLAuthMethodsTool := tools.GetACLAuthMethodsTool(logger)
	hcServer.AddTool(getACLAuthMethodsTool.Tool, getACLAuthMethodsTool.Handler)

	getACLAuthMethodTool := tools.GetACLAuthMethodTool(logger)
	hcServer.AddTool(getACLAuthMethodTool.Tool, getACLAuthMethodTool.Handler)

	getACLBindingRulesTool := tools.GetACLBindingRulesTool(logger)
	hcServer.AddTool(getACLBindingRulesTool.Tool, getACLBindingRulesTool.Handler)

	getACLBindingRuleTool := tools.GetACLBindingRuleTool(logger)
	hcServer.AddTool(getACLBindingRuleTool.Tool, getACLBindingRuleTool.Handler)

	identityTool := tools.GetIdentity(logger)
	hcServer.AddTool(identityTool.Tool, identityTool.Handler)
}
