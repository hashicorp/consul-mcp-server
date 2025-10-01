# <img src="public/images/consul-fill-color-24.svg" width="30" align="left" style="margin-right: 12px; margin-top: 5px"/> Consul MCP Server

The Consul MCP Server is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction)
server that provides seamless integration with HashiCorp Consul APIs, enabling advanced
automation and interaction capabilities for service discovery, configuration management, and service mesh operations.

## Features

- **Dual Transport Support**: Both Stdio and StreamableHTTP transports
- **Service Discovery**: Query and manage services, nodes, and health checks in the Consul catalog
- **Key-Value Store**: Access and manage Consul's distributed key-value store
- **Service Mesh**: Interact with Consul Connect for service mesh functionality including intentions and certificates
- **ACL Management**: Manage Access Control Lists including tokens, policies, roles, and auth methods
- **Agent Operations**: Monitor and configure Consul agents, including health checks and services
- **Cluster Operations**: Access operator tools for cluster management, autopilot, and Raft operations
- **Container Ready**: Docker support for easy deployment

> **Caution:** The outputs and recommendations provided by the MCP server are generated dynamically and may vary based on the query, model, and the connected MCP server. Users should **thoroughly review all outputs/recommendations** to ensure they align with their organization's **security best practices**, **compliance requirements**, and **Consul deployment policies** before implementation.

> **Security Note:** When using the StreamableHTTP transport in production, always configure the `MCP_ALLOWED_ORIGINS` environment variable to restrict access to trusted origins only. This helps prevent DNS rebinding attacks and other cross-origin vulnerabilities.

## Prerequisites

1. A running Consul cluster or agent that the MCP server can connect to
2. For containerized deployment, [Docker](https://www.docker.com/) installed and running
3. Appropriate Consul ACL tokens if ACL is enabled in your Consul cluster

## Transport Support

The Consul MCP Server supports multiple transport protocols:

### 1. Stdio Transport (Default)
Standard input/output communication using JSON-RPC messages. Ideal for local development and direct integration with MCP clients.

### 2. StreamableHTTP Transport
Modern HTTP-based transport supporting both direct HTTP requests and Server-Sent Events (SSE) streams. This is the recommended transport for remote/distributed setups.

**Features:**
- **Endpoint**: `http://{hostname}:8080/mcp`
- **Health Check**: `http://{hostname}:8080/health`
- **Environment Configuration**: Set `TRANSPORT_MODE=streamable-http` or `TRANSPORT_PORT=8080` to enable

**Environment Variables:**

| Variable | Description                                                                                                                                   | Default                 |
|----------|-----------------------------------------------------------------------------------------------------------------------------------------------|-------------------------|
| `TRANSPORT_MODE` | Set to `streamable-http` to enable HTTP transport (legacy `http` value still supported)                                                       | `stdio`                 |
| `TRANSPORT_HOST` | Host to bind the HTTP server                                                                                                                  | `127.0.0.1`             |
| `TRANSPORT_PORT` | HTTP server port                                                                                                                              | `8080`                  |
| `MCP_ENDPOINT` | HTTP server endpoint path                                                                                                                     | `/mcp`                  |
| `MCP_SESSION_MODE` | Session mode: `stateful` or `stateless`                                                                                                       | `stateful`              |
| `MCP_ALLOWED_ORIGINS` | Comma-separated list of allowed origins for CORS                                                                                              | `""` (empty)            |
| `MCP_CORS_MODE` | CORS mode: `strict`, `development`, or `disabled`                                                                                             | `strict`                |
| `CONSUL_HTTP_ADDR` | Consul agent HTTP API address                                                                                                                 | `http://127.0.0.1:8500` |
| `CONSUL_HTTP_TOKEN` | Consul ACL token for authentication                                                                                                           | `""` (empty)            |
| `CONSUL_SKIP_VERIFY` | Skip TLS certificate verification (use only for development/testing)                                                                          | `false`                 |
| `CONSUL_ENTERPRISE` | Enable Consul Enterprise features and API endpoints                                                                                           | `true`                   |
| `CONSUL_MCP_SERVER_READ_GITHUB_RESOURCES` | For latest Consul context, the flag enables the fetching resource from https://github.com/hashicorp/consul/blob/main/website/content/api-docs | `true`                  |

## TLS Configuration

The Consul MCP Server supports secure connections to Consul clusters with proper TLS configuration. Here are the key settings:

### Certificate Verification

By default, the MCP server verifies TLS certificates when connecting to HTTPS Consul endpoints. If your Consul cluster uses self-signed certificates or certificates that cannot be verified against the system's certificate store, you may encounter errors like:

```
tls: failed to verify certificate: x509: "________" certificate is not trusted
```

### Disabling Certificate Verification (Development/Testing Only)

**⚠️ Security Warning**: Only disable certificate verification in development or testing environments. Never use this in production.

To disable TLS certificate verification, set the environment variable:

```bash
export CONSUL_SKIP_VERIFY=true
```

### Example Configurations for Development

**Development with self-signed certificates:**
```bash
export CONSUL_HTTP_ADDR=https://consul.example.com:8501
export CONSUL_HTTP_TOKEN=your-acl-token
export CONSUL_SKIP_VERIFY=true
```

**Production with proper certificates:**
```bash
export CONSUL_HTTP_ADDR=https://consul.example.com:8501
export CONSUL_HTTP_TOKEN=your-acl-token
export CONSUL_SKIP_VERIFY=false
# CONSUL_SKIP_VERIFY should remain false (default)
```

**Docker example with TLS skip verification:**
```bash
docker run -i --rm \
  -e CONSUL_HTTP_ADDR=https://host.docker.internal:8501 \
  -e CONSUL_HTTP_TOKEN=your-token \
  -e CONSUL_SKIP_VERIFY=true \
  hashicorp/consul-mcp-server
```

## Command Line Options

```bash
# Stdio mode
consul-mcp-server stdio [--log-file /path/to/log]

# StreamableHTTP mode
consul-mcp-server streamable-http [--transport-port 8080] [--transport-host 127.0.0.1] [--mcp-endpoint /mcp] [--log-file /path/to/log]
```

## Session Modes

The Consul MCP Server supports two session modes when using the StreamableHTTP transport:

- **Stateful Mode (Default)**: Maintains session state between requests, enabling context-aware operations.
- **Stateless Mode**: Each request is processed independently without maintaining session state, which can be useful for high-availability deployments or when using load balancers.

To enable stateless mode, set the environment variable:
```bash
export MCP_SESSION_MODE=stateless
```

## Installation

### Usage with VS Code

Add the following JSON block to your User Settings (JSON) file in VS Code. You can do this by pressing `Ctrl + Shift + P` and typing `Preferences: Open User Settings (JSON)`. 

More about using MCP server tools in VS Code's [agent mode documentation](https://code.visualstudio.com/docs/copilot/chat/mcp-servers).

```json
{
  "mcp": {
    "servers": {
      "consul": {
        "command": "docker",
        "args": [
          "run",
          "-i",
          "--rm",
          "-e", "CONSUL_HTTP_ADDR=http://host.docker.internal:8500",
          "hashicorp/consul-mcp-server"
        ]
      }
    }
  }
}
```

Optionally, you can add a similar example (i.e. without the mcp key) to a file called `.vscode/mcp.json` in your workspace. This will allow you to share the configuration with others.

```json
{
  "servers": {
    "consul": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e", "CONSUL_HTTP_ADDR=http://host.docker.internal:8500",
        "hashicorp/consul-mcp-server"
      ]
    }
  }
}
```

Copilot client mode

```json
{
  servers: {
    "consul-cluster-dc1": {
      "url": "http://<mcp-server-dc1-address>:<port>",
      "headers": {
        "X-Consul-Address": "<consul-address-dc1>",
        "X-Consul-Token": "<consul-acl-token-dc1>"
      }
    },
    "consul-cluster-dc2": {
      "url": "http://<mcp-server-dc2-address>:<port>",
      "headers": {
        "X-Consul-Address": "<consul-address-dc2>",
        "X-Consul-Token": "<consul-acl-token-dc2>"
      }
    }
  },
  "inputs": []
}
```

### Usage with Claude Desktop / Amazon Q Developer / Amazon Q CLI

More about using MCP server tools in Claude Desktop [user documentation](https://modelcontextprotocol.io/quickstart/user).
Read more about using MCP server in Amazon Q from the [documentation](https://docs.aws.amazon.com/amazonq/latest/qdeveloper-ug/qdev-mcp.html).

```json
{
  "mcpServers": {
    "consul-dc1": {
      "command": "consul-mcp-server",
      "args": ["--host", "localhost", "--port", "8080"],
      "env": {
        "CONSUL_HTTP_ADDR": "https://consul-dc1.example.com:8501",
        "CONSUL_HTTP_TOKEN": "${CONSUL_DC1_TOKEN}"
      }
    },
    "consul-dc2": {
      "command": "consul-mcp-server",
      "args": ["--host", "localhost", "--port", "8081"],
      "env": {
        "CONSUL_HTTP_ADDR": "https://consul-dc2.example.com:8501",
        "CONSUL_HTTP_TOKEN": "${CONSUL_DC2_TOKEN}"
      }
    }
  }
}
```

## Tool Configuration

### Available Toolsets

The following sets of tools are available for interacting with Consul:

| Toolset | Tools | Description |
|---------|-------|-------------|
| `catalog` | `get_catalog_services`, `get_catalog_nodes`, `get_catalog_service`, `get_catalog_connect`, `get_catalog_node`, `get_catalog_datacenters`, `get_catalog_gateway_services` | Query and explore services and nodes in the Consul catalog |
| `agent` | `get_agent_self`, `get_agent_config`, `get_agent_members`, `get_agent_metrics`, `get_agent_host`, `get_agent_version`, `get_agent_reload` | Monitor and configure Consul agents |
| `health` | `get_health_node`, `get_health_checks`, `get_health_service`, `get_health_connect`, `get_health_ingress`, `get_health_state` | Query health information for services and nodes |
| `kv` | `get_kv`, `get_kv_keys`, `get_kv_recursive` | Access and manage the Consul key-value store |
| `acl` | `get_acl_tokens`, `get_acl_policies`, `get_acl_roles`, `get_acl_auth_methods`, `get_acl_binding_rules`, `get_acl_templated_policies` | Manage Access Control Lists and authentication |
| `connect` | `get_connect_ca_roots`, `get_connect_ca_configuration`, `get_connect_intentions`, `get_connect_intention`, `get_connect_intention_match`, `get_connect_intention_check` | Manage Consul Connect service mesh features |
| `operator` | `get_operator_autopilot_*`, `get_operator_keyring`, `get_operator_license`, `get_operator_raft_*`, `get_operator_usage` | Access cluster operational tools and configuration |
| `session` | `get_session`, `get_session_node`, `get_session_list` | Manage Consul sessions for distributed locking |
| `status` | `get_status_leader`, `get_status_peers` | Query cluster status and leadership information |
| `peering` | `get_peerings`, `get_peering`, `get_peering_exported_services` | Manage cluster peering relationships |
| `config` | `get_config_entries`, `get_config_entry` | Access Consul configuration entries |
| `discovery` | `get_discovery_chain` | Query service discovery chains |
| `query` | `get_query`, `get_query_by_id`, `get_query_execute`, `get_query_explain` | Execute and manage prepared queries |
| `namespaces` | `get_namespaces`, `get_namespace` | Manage Consul Enterprise namespaces (Enterprise only) |
| `identity` | Various identity-related tools | Manage service identity and certificates |

## Resource Configuration

### Available Resources

| Resource URI | Description |
|--------------|-------------|
| `consul://connect/ca/roots` | Consul Cluster Identity - Provides access to the Consul Connect CA root certificates and cluster identity information |
| `consul://api-docs/*` | Consul API Documentation - Dynamic access to official Consul API documentation from the GitHub repository |

### Install from source

Use the latest release version:

```console
go install github.com/hashicorp/consul-mcp-server/cmd/consul-mcp-server@latest
```

Use the main branch:

```console
go install github.com/hashicorp/consul-mcp-server/cmd/consul-mcp-server@main
```

```json
{
  "mcp": {
    "servers": {
      "consul": {
        "command": "/path/to/consul-mcp-server",
        "args": ["stdio"]
      }
    }
  }
}
```
