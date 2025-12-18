package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/a2aproject/a2a-go/a2a"
	"github.com/a2aproject/a2a-go/a2aclient"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/remoteagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/web"
	a2aweb "google.golang.org/adk/cmd/launcher/web/a2a"
	"google.golang.org/adk/cmd/launcher/web/webui"
	"google.golang.org/adk/cmd/launcher/web/api"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
 	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/api/idtoken"
	"google.golang.org/genai"
)

// fetchAgentCard retrieves the agent card available at the provided host url.
func fetchAgentCard(hostUrl string) (a2a.AgentCard, error) {
	
	ctx := context.Background()

	// Create a new HTTP client
	client, err := idtoken.NewClient(ctx, hostUrl)
	if err != nil {
		return a2a.AgentCard{}, fmt.Errorf("idtoken.NewClient: %w", err)
	}

	// Prepare the HTTP request
	httpReq, err := http.NewRequest("GET", hostUrl+"/.well-known/agent-card.json", nil)
	if err != nil {
		return a2a.AgentCard{}, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return a2a.AgentCard{}, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	// Parse response body
	httpRespBytes, err := io.ReadAll(httpResp.Body)
    if err != nil {
        return a2a.AgentCard{}, fmt.Errorf("failed to read response body: %v", err)
    }

	// Check for error status codes
	if httpResp.StatusCode < 200 || httpResp.StatusCode > 300 {
		httpResp.Body.Close()
		return a2a.AgentCard{}, fmt.Errorf("unexpected status code %d: %s", httpResp.StatusCode, string(httpRespBytes))
	}

	// Parse response body into expected type
	agentCard := a2a.AgentCard{}
	err = json.Unmarshal(httpRespBytes, &agentCard)
	if err != nil {
		return a2a.AgentCard{}, err
	}

	return agentCard, nil
}

// newRemoteAgent registers a new remote A2A agent through its AgentCard
func newRemoteAgent(hostUrl string) (agent.Agent, error) {

	// Fetch the agent card
	agentCard, err := fetchAgentCard(hostUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent card: %v", err)
	}

	// Create new AuthInterceptor to extend the default client factory
	authInterceptor, err := NewAuthInterceptor(hostUrl)
	if err != nil {
		return nil, err
	}

	// Initialise remote agent
	remoteAgent, err := remoteagent.NewA2A(remoteagent.A2AConfig{
		Name:        agentCard.Name,
		Description: agentCard.Description,
		ClientFactory: a2aclient.NewFactory(
			a2aclient.WithInterceptors(authInterceptor),
		),
		AgentCard: &agentCard,
	})
	if err != nil {
		return nil, err
	}
	return remoteAgent, nil
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

    // Initialise new model
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    	Backend: genai.BackendVertexAI,
    	Project: os.Getenv("GOOGLE_CLOUD_PROJECT"),
    	Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
    })
    if err != nil {
        log.Fatalf("Failed to create model: %v", err)
    }

	// Register remote agents and then convert to tools
    var remoteAgentTools []tool.Tool
	urls := []string{
		// TODO: Add remote agent Cloud Run host URLs
		// Format: https://{service-name}-{projectNumber}.{location}.run.app 
	}
    for _, url := range urls {
		remoteAgent, err := newRemoteAgent(url)
		if err != nil {
			log.Fatalf("failed to create remote agent: %v", err)
		}
		remoteAgentTool := agenttool.New(remoteAgent, &agenttool.Config{
			SkipSummarization: false,
		})
		remoteAgentTools = append(remoteAgentTools, remoteAgentTool)
    }

	// Create new agent
    a, err := llmagent.New(llmagent.Config{
		Name: "Orchestrator Agent",
		Model: model,
		Description: "Orchestrator agent to assist with user queries and directing them to specialist agents",
		Instruction: "You are an orchesrator agent. Make use of the available remote agents to execute tasks as required.",
		Tools: remoteAgentTools,
	})
    if err != nil {
        log.Fatalf("Failed to create agent: %v", err)
    }

    config := &launcher.Config{
        AgentLoader: agent.NewSingleLoader(a),
		SessionService: session.InMemoryService(),
    }

	// Set web launcher to interact via the ADK webui
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

