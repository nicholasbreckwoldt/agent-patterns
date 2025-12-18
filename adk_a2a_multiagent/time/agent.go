package main

import (
	"context"
	"log"
    "os"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
    "google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/cmd/launcher/web/a2a"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
    "google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

// GetCurrentTimeParams defines the arguments for the getCurrentTime tool.
type GetCurrentTimeParams struct {
    // The name of the city
    City string `json:"city"`
}

// GetCurrentTimeOutput defines the output of the getCurrentTime tool.
type GetCurrentTimeOutput struct {
    // The name of the city
    City string `json:"city"`
    // The time
    Time string `json:"time"`
}

// Tool handler
func getCurrentTimeHandler(ctx tool.Context, args GetCurrentTimeParams) (GetCurrentTimeOutput, error) {
    return GetCurrentTimeOutput{
        City:    args.City,
        Time:    "10:30 AM",
    }, nil
}

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

    // Define time tool
    getCurrentTimeTool, err := functiontool.New(
        functiontool.Config{
            Name:        "get_current_time",
            Description: "Returns the current time in a specified city.",
        },
        getCurrentTimeHandler,
    )
    if err != nil {
        log.Fatalf("Failed to create tool: %v", err)
    }

    // Define LLMAgent
    a, err := llmagent.New(llmagent.Config{
        Name:        "time_agent",
        Model:       model,
        Description: "Assists with time related queries",
        Instruction: "You are a helpful assistant that can provide the time information across different world cities.",
        Tools: []tool.Tool{
            getCurrentTimeTool,
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

	webLauncher := web.NewLauncher(a2a.NewLauncher())
	_, err = webLauncher.Parse([]string{
        // TODO: Populate CLOUD_RUN_HOST_URL before deploying
        // Format: https://{service-name}-{projectNumber}.{location}.run.app 
		"--port", "8080", "a2a", "--a2a_agent_url", "{CLOUD_RUN_HOST_URL}",
	})
	if err != nil {
		log.Fatalf("webLauncher.Parse() error = %v", err)
	}

    if err := webLauncher.Run(context.Background(), config); err != nil {
		log.Fatalf("webLauncher.Run() error = %v", err)
	}
}

