# Multi-agent using ADK and A2A

## Overview 
This project demonstrates a multi-agent system built with the Google Agent Development Kit (ADK). It implement the "agent-as-a-tool" pattern consisting of an orchestrator agent and two specialist agents (time and weather agents) communicating via A2A.

- **Orchestrator Agent**: This is the main entry point for user interaction.

- **Time Agent**: A specialist agent that provides the current time for a given city.

- **Weather Agent**: A specialist agent that provides the current weather report for a given city.

## Before you begin

This project assumes you have the following:

- The `gcloud` CLI installed and authenticated.
- A Google Cloud Project.
- A service account with the required roles and permissions (i.e. `roles/run.invoker`, etc)
- A service account key file for obtaining local credentials.
- GoLang installation

## Running the project (Makefile)

The `Makefile` provides a convenient way to build, deploy, and run the agent system.

### Available Commands

- `make deploy_time_agent`: Builds the Time Agent from source and deploys it to Cloud Run.
- `make deploy_weather_agent`: Builds the Weather Agent from source and deploys it to Cloud Run.
- `make run_orchestrator`: Runs the orchestrator agent locally.
- `make all`: Builds and deploys both the time and weather agents to Cloud Run, then runs the orchestrator agent locally.

### Example Usage

To deploy both remote agents and run the local Orchestrator via the ADK WebUI, run the following commmand in the `adk_a2a_multiagent` folder:

```bash
make all
```

Once the orchestrator is running, you can access the web UI at [http://localhost:8080/ui](http://localhost:8080/ui).

**Important Note on Agent URLs:**
The `deploy_time_agent` and `deploy_weather_agent` commands deploy the respective agents to Google Cloud Run. You will need to update the URLs within the `orchestrator/agent.go` file to match the deployed service URLs for the orchestrator to correctly discover and interact with them.