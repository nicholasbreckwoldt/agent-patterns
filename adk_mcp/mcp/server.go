package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"fmt"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Args for the GetCurrentTime tool
type GetCurrentTimeArgs struct {
	City string `json:"city"`
}

// GetCurrentTime retrieves the current time for a given city
func GetCurrentTime(ctx context.Context, req *mcp.CallToolRequest, args GetCurrentTimeArgs) (*mcp.CallToolResult, any, error) {

	// Mock city time data
	times := map[string]string{
		"london": "10.30am",
		"new york": "5.30am",
	}
	currentTime, ok := times[strings.ToLower(args.City)]
	if !ok {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Unfortunately the time for %s is currently unavailable", args.City)},
			},
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("The current time is %s", currentTime)},
		},
	}, nil, nil
}

// createMCPServer creates an MCP server with registered tools
func createMCPServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "My Remote MCP Server",
		Version: "1.0.0",
	}, nil)
	
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_current_time",
		Description: "Retrieves the current time",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"city": {
					Description: "City for which to retrieve the current time.",
					Type:        "string",
				},
			},
			Required: []string{"city"},
		},
	}, GetCurrentTime)

	return server
}

func main() {
	// Create the MCP server.
	mcpServer := createMCPServer()

	// Create HTTP handler
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return mcpServer
	}, &mcp.StreamableHTTPOptions{})

	// Listen and serve the MCP Server
	if err := http.ListenAndServe(":8080", handler); err != nil {
        log.Fatal(err)
    }
}