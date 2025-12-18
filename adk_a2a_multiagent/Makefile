# Variables
GOOGLE_CLOUD_PROJECT=google-project
GOOGLE_CLOUD_LOCATION=google=region
SERVICE_ACCOUNT=service-account
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account/key.json

export GOOGLE_CLOUD_PROJECT
export GOOGLE_CLOUD_LOCATION
export GOOGLE_APPLICATION_CREDENTIALS
export GOOGLE_GENAI_USE_VERTEXAI=TRUE

# Default target
all: deploy_weather_agent deploy_time_agent run_orchestrator

# Helper target
help:
	@echo "Available commands:"
	@echo "  make deploy_time_agent    - Builds Time Agent from source and deploys to Cloud Run"
	@echo "  make deploy_weather_agent - Builds Weather Agent from source and deploys to Cloud Run"
	@echo "  make all                  - Builds and deploys both agents from source"

# Build and deploy the time agent target
deploy_time_agent:
	@echo "Deploying `Time Agent` to Cloud Run..."
	gcloud run deploy time-agent \
		--source=./time \
		--region=$(GOOGLE_CLOUD_LOCATION) \
		--service-account=$(SERVICE_ACCOUNT) \
		--project $(GOOGLE_CLOUD_PROJECT) \
		--set-env-vars=GOOGLE_CLOUD_PROJECT=$(GOOGLE_CLOUD_PROJECT) \
		--set-env-vars=GOOGLE_GENAI_USE_VERTEXAI=TRUE \
		--set-env-vars=GOOGLE_CLOUD_LOCATION=$(GOOGLE_CLOUD_LOCATION) \
		--no-allow-unauthenticated
	@echo "Time agent deployment complete!"

# Build and deploy the weather agent target
deploy_weather_agent:
	@echo "Deploying `Weather Agent` to Cloud Run..."
	gcloud run deploy weather-agent \
		--source=./weather \
		--region $(GOOGLE_CLOUD_LOCATION) \
		--project $(GOOGLE_CLOUD_PROJECT) \
		--service-account=$(SERVICE_ACCOUNT) \
		--set-env-vars=GOOGLE_CLOUD_PROJECT=$(GOOGLE_CLOUD_PROJECT) \
		--set-env-vars=GOOGLE_GENAI_USE_VERTEXAI=TRUE \
		--set-env-vars=GOOGLE_CLOUD_LOCATION=$(GOOGLE_CLOUD_LOCATION) \
		--no-allow-unauthenticated
	@echo "Weather agent deployment complete!"

run_orchestrator:
	cd ./orchestrator && \
	go mod tidy && \
	go run .
