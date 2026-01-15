# Variables
GOOGLE_CLOUD_PROJECT=google-project
SERVICE_ACCOUNT=service-account
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account/key.json
GOOGLE_CLOUD_LOCATION=google-region

export GOOGLE_CLOUD_PROJECT
export GOOGLE_CLOUD_LOCATION

# Build and deploy the MCP server
deploy_mcp:
	@echo "Deploying MCP server to Cloud Run..."
	gcloud run deploy mcp-server \
		--source=./mcp \
		--service-account=$(SERVICE_ACCOUNT) \
		--project $(GOOGLE_CLOUD_PROJECT) \
		--region=$(GOOGLE_CLOUD_LOCATION) \
		--no-allow-unauthenticated
	@echo "MCP server deployment complete!"

# Run the ADK agent
run_agent:
	cd ./agent && \
	go mod tidy && \
	go run .
