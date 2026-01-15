package main

import (
	"context"
	"log"
    "os"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
    "google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

    a2aweb "google.golang.org/adk/cmd/launcher/web/a2a"
	"google.golang.org/adk/cmd/launcher/web/webui"
	"google.golang.org/adk/cmd/launcher/web/api"
)

func init() {
    if os.Getenv("GOOGLE_CLOUD_PROJECT") == "" {
        log.Fatalf("GOOGLE_CLOUD_PROJECT env not set")
    }
      if os.Getenv("GOOGLE_CLOUD_LOCATION") == "" {
        log.Fatalf("GOOGLE_CLOUD_LOCATION env not set")
    }
}

func main() {
    ctx := context.Background()

    // Create model
    model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    	Backend: genai.BackendVertexAI,
    	Project: os.Getenv("GOOGLE_CLOUD_PROJECT"),
    	Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
    })
    if err != nil {
        log.Fatalf("Failed to create model: %v", err)
    }

    // Create the MCP Toolset
    mcpToolSet, err := generateStreamableHTTPMCPToolSet(ctx, 
        "YOUR_CLOUD_RUN_HOST", "/mcp") // TODO: Replace with your CloudRun host URL
    if err != nil {
		log.Fatalf("Failed to create MCP ToolSet: %v", err)
	}

    // Define LLMAgent
    a, err := llmagent.New(llmagent.Config{
        Name:        "time_agent",
        Model:       model,
        Description: "Assists with time related queries",
        Instruction: "You are a helpful assistant that can use available tools to help user",
        Toolsets: []tool.Toolset{
            mcpToolSet,
        },
    })
    if err != nil {
        log.Fatalf("Failed to create agent: %v", err)
    }

    // Launch and serve
    config := &launcher.Config{
        AgentLoader: agent.NewSingleLoader(a),
        SessionService: session.InMemoryService(),
    }

	webLauncher := web.NewLauncher(
        a2aweb.NewLauncher(),
		webui.NewLauncher(),
		api.NewLauncher(),
    )
    _, err = webLauncher.Parse([]string{"--port", "8080", "webui", "api"})
	if err != nil {
		log.Fatalf("webLauncher.Parse() error = %v", err)
	}

    if err := webLauncher.Run(context.Background(), config); err != nil {
		log.Fatalf("webLauncher.Run() error = %v", err)
	}
}

