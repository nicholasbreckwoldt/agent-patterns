# MCP + ADK

## Overview 
This project demonstrates a simple Model Context Protocol (MCP) server built
using the official GoLang MCP SDK. The deployed MCP Server is then connected to an
ADK agent which can then invoke the available tools via MCP.


- **mcp**: This folder contains the simple MCP Server implementation.

- **agent**: This folder contains a basic ADK agent that invokes the MCP server.

## Before you begin

This project assumes you have the following:

- The `gcloud` CLI installed and authenticated.
- A Google Cloud Project.
- A service account with the required roles and permissions (i.e. `roles/run.invoker`, etc)
- A service account key file for obtaining local credentials.
- GoLang installation

## Running the project (Makefile)

The `Makefile` provides a convenient way to deploy the MCP Server and run an ADK agent with the provided tools.

### Available Commands

- `make deploy_mcp`: Builds the MCP Server and deploys it to Cloud Run.
- `make run_agent`: Runs the agent locally.

### Example Usage

To build and deploy the MCP Server, run the following commmand in the `adk_mcp` folder:

```bash
make deploy_mcp
```
To launch the simple ADK agent, run the following:

```bash
make run_agent
```

Once the agent is running, you can access the web UI at [http://localhost:8080/ui](http://localhost:8080/ui).