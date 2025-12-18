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

// GetWeatherParams defines the arguments for the getWeather tool.
type GetWeatherParams struct {
    // The name of the city
    City string `json:"city"`
}

// GetWeatherOutput defines the output of the getWeather tool.
type GetWeatherOutput struct {
    // Weather report
    Report string `json:"report"`
}

// Tool handler
func getWeatherHandler(ctx tool.Context, args GetWeatherParams) (GetWeatherOutput, error) {
    return GetWeatherOutput{
        Report: "The temperature is 20 degrees celisus with no chance of rain",
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

    // Define weather tool
    getWeatherTool, err := functiontool.New(
        functiontool.Config{
            Name:        "get_weather",
            Description: "Returns the current weather report in a specified city.",
        },
        getWeatherHandler,
    )
    if err != nil {
        log.Fatalf("Failed to create tool: %v", err)
    }

    // Define LLMAgent
    a, err := llmagent.New(llmagent.Config{
        Name:        "weather_agent",
        Model:       model,
        Description: "Assists with weather related queries",
        Instruction: "You are a helpful assistant can assist users with weather related queries",
        Tools: []tool.Tool{
            getWeatherTool,
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
