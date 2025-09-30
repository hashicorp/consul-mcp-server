package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul-mcp-server/pkg/client"
	"github.com/hashicorp/consul-mcp-server/pkg/utils"
	"github.com/hashicorp/consul/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Anonymous access could cause rate limiting, so we point to the raw files in the main branch, hence caching the URL
const ConsulGuideRawURL = "https://github.com/hashicorp/consul/blob/main/website/content/api-docs"

var api_docs = map[string]ApiDoc{}

type ApiDoc struct {
	Name string
	URL  string
}

type GitHubAPIResponse struct {
	Payload struct {
		Tree struct {
			Items []GitHubFileInfo `json:"items"`
		} `json:"tree"`
	} `json:"payload"`
}

type GitHubFileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	ContentType string `json:"contentType"`
}

// RegisterResources adds the new resource
func RegisterResources(hcServer *server.MCPServer, logger *log.Logger) {
	// Called from main.go to register resources
	ctx := context.Background()
	// Get a simple http client to access the GitHub raw files
	httpClient := client.NewHttpClientFromContext(ctx, logger)

	// Add the Consul identity resource
	hcServer.AddResource(consulidentityResource(logger))

	// Reading could be disabled via env var for testing with local files
	// e.g. export CONSUL_MCP_SERVER_READ_GITHUB_RESOURCES=false
	readGithubResource := func() bool {
		if val := os.Getenv("CONSUL_MCP_SERVER_READ_GITHUB_RESOURCES"); val != "" {
			return strings.ToLower(val) != "false"
		}
		return true
	}

	if len(api_docs) == 0 && readGithubResource() {
		// not initialized yet, so fetch
		fetchedDocs, err := fetchConsulAPIDocsFiles(ctx, ConsulGuideRawURL, httpClient)
		if err != nil {
			logger.Errorf("Failed to fetch Consul API docs from Consul: %v", err)
			return
		}
		// Append to existing api_docs map rather than replacing
		for key, value := range fetchedDocs {
			api_docs[key] = value
		}
	}

	// fetch each doc and register as resource
	for _, doc := range api_docs {
		hcServer.AddResource(consulAPIDocResource(doc, logger))
	}
}

// Provides the agent the identity of the Consul server
func consulidentityResource(logger *log.Logger) (mcp.Resource, server.ResourceHandlerFunc) {
	resourceURI := "consul://connect/ca/roots"
	description := fmt.Sprintf("Consul Cluster Identity")

	return mcp.NewResource(
			resourceURI,
			description,
			mcp.WithMIMEType("text/plain"),
			mcp.WithResourceDescription(description),
		),
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			consulClient, err := client.GetGetConsulHttpClientFromContext(ctx, logger)
			if err != nil {
				return nil, utils.LogAndReturnError(logger, "getting http client for Consul", err)
			}

			// Build query parameters
			queryParams := url.Values{}

			rootsResp, err := consulClient.Get("connect/ca/roots", queryParams)
			if err != nil {
				return nil, utils.LogAndReturnError(logger, "fetching Connect CA roots from consul", err)
			}

			// Create types.ConnectCARoots from rootsResp
			var caRoots api.CARootList
			if err := json.Unmarshal(rootsResp, &caRoots); err != nil {
				return nil, utils.LogAndReturnError(logger, "failed to unmarshal Connect CA roots response", err)
			}

			// convert rootsResp i.e. bytes[] to text
			rootsText := strings.TrimSpace(caRoots.TrustDomain)
			var contents []mcp.ResourceContents
			contents = append(contents, mcp.TextResourceContents{
				MIMEType: "text/plain",
				URI:      request.Params.URI,
				Text:     fmt.Sprintf("Consul Cluster Identity: %s", rootsText),
			})
			return contents, nil
		}
}

func consulAPIDocResource(doc ApiDoc, logger *log.Logger) (mcp.Resource, server.ResourceHandlerFunc) {
	resourceURI := doc.URL
	description := fmt.Sprintf("Consul API Documentation: %s", doc.Name)

	return mcp.NewResource(
			resourceURI,
			description,
			mcp.WithMIMEType("text/markdown"),
			mcp.WithResourceDescription(description),
		),
		func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			httpClient := client.NewHttpClientFromContext(ctx, logger)

			resp, err := httpClient.Get(resourceURI)
			if err != nil {
				return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching %s markdown", doc.Name), err)
			}
			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching %s markdown, status not ok", doc.Name), fmt.Errorf("status: %s", resp.Status))
			}
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				// Ignore this file and continue
				return nil, utils.LogAndReturnError(logger, fmt.Sprintf("fetching %s markdown", doc.Name), err)
			}

			bodyStr := string(body)
			logger.Infof("Fetched %s markdown, length %d", doc.Name, len(bodyStr))
			var contents []mcp.ResourceContents
			contents = append(contents, mcp.TextResourceContents{
				MIMEType: "text/markdown",
				URI:      doc.URL,
				Text:     bodyStr,
			})

			return contents, nil
		}
}

func fetchConsulAPIDocsFiles(ctx context.Context, apiURL string, httpClient *http.Client) (map[string]ApiDoc, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch API docs files: status %d", resp.StatusCode)
	}

	var apiResp GitHubAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub API response: %w", err)
	}

	result := make(map[string]ApiDoc)
	for _, item := range apiResp.Payload.Tree.Items {
		// Only process files (not directories)
		if item.ContentType == "file" && (strings.HasSuffix(item.Name, ".mdx") || strings.HasSuffix(item.Name, ".md")) {
			// Extract the base name without extension for the key
			name := strings.TrimSuffix(item.Name, ".mdx")
			name = strings.TrimSuffix(name, ".md")

			// Construct the raw GitHub URL for the file
			downloadURL := fmt.Sprintf("%s/%s", apiURL, item.Name)
			result[name] = ApiDoc{Name: name, URL: downloadURL}
		} else if item.ContentType == "directory" {
			// recursively fetch files in the directory
			dirUrl := fmt.Sprintf("%s/%s", apiURL, item.Name)
			subFiles, err := fetchConsulAPIDocsFiles(ctx, dirUrl, httpClient)
			if err != nil {
				return nil, err
			}
			// Merge subFiles into result
			for k, v := range subFiles {
				result[k] = v
			}
		}
	}

	return result, nil
}
