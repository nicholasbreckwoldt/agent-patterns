package main

import (
	"context"
	"google.golang.org/adk/tool"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/tool/mcptoolset"
	"google.golang.org/api/idtoken"
)

// generateStreamableHTTPMCPToolSet creates a StreamableHTTP MCPToolSet connection.
func generateStreamableHTTPMCPToolSet(ctx context.Context, host, endpoint string) (tool.Toolset, error) {

	// Establish client for transport
	httpClient, err := idtoken.NewClient(ctx, host)
	if err != nil {
		return nil, err
	}

	// Prepare MCP ToolSet
	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: &mcp.StreamableClientTransport{
			Endpoint: host+endpoint,
			HTTPClient: httpClient,
    	},
	})
	if err != nil {
		return nil, err
	}
	return mcpToolSet, nil
}