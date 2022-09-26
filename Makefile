PROJECT_NAME:=latency-aggregator
TARGET:=./cmd/latency-aggregator/main.go
APP_INTERNAL_DEPTH := $(shell find ./internal -type d -printf '%d\n' | sort -rn | head -1)
APP_VERSION=0.0.1

include Makefile.tests.mk
include Makefile.coverage.mk
include Makefile.docker.mk

.PHONY: all
all: help

.PHONY: help
help: ## Show commands and description
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: ## Install all dependencies
	@echo "Get the dependencies..."
	@make dep --silent 

	@echo "Install staticcheck to lint..."
	@go install honnef.co/go/tools/cmd/staticcheck@2022.1.3

	@echo "Install gosec..."
	go install github.com/securego/gosec/v2/cmd/gosec@v2.13.1
	
	@echo "Configuring hooks..."
	@git config core.hooksPath hooks/
	@chmod +x ./hooks/pre-commit

.PHONY: run-docker
run-docker: ## Run web application and dependencies inside container
	@docker-compose -f docker-compose.app.yml up -d --build 
	@docker-compose -f docker-compose.yml up -d
	
.PHONY: run
run: dep-dev-run ## Run web application and dependencies
	@ENVIRONMENT=development APP_VERSION="v-$(shell git rev-parse --short HEAD)" go run ${TARGET}

.PHONY: down
down: dep-dev-stop ## Run application and dependencies
	@docker-compose down -v

.PHONY: lint
lint: ## Lint the files
	@staticcheck ./... # TODO: Dosn't work with go 1.18

.PHONY: vet
vet:  ## Run go vet in project
	@go vet ./...

.PHONY: gosec
gosec: ## Run gosec in project
	@gosec ./...

.PHONY: dep
dep: ## Get the dependencies
	@go get -v -d ./...

.PHONY: build
build: dep ## Build the binary file
	env GOOS=linux GOARCH=amd64 go build -v -o bin/${PROJECT_NAME} ${TARGET}

.PHONY: clean
clean: ## Remove previous build
	@rm -f bin/$(PROJECT_NAME)

.PHONY: scan-code
scan-code:
	docker run --rm \
		-e HORUSEC_CLI_FILES_OR_PATHS_TO_IGNORE="*tmp*, **/.vscode/**, **/docs/**, **/node_modules/**, **/.horusec/**, **/.trivy/**" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD):/src horuszup/horusec-cli:latest \
		horusec start -p /src -P $(PWD)

.PHONY: scan-image
scan-image:
	docker run --rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v ${PWD}/.trivy/.cache:/root/.cache/ \
		aquasec/trivy:0.18.3 \
		incident-webhook_incident-webhook

.PHONY: scan
scan: lint vet gosec scan-code scan-image

.PHONY: go-to-uml
go-to-uml:
	goplantuml -recursive .  > "docs/diagram_file_name.puml"
