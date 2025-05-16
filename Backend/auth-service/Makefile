include .env.example
export

LOCAL_BIN:=$(CURDIR)/bin
BASE_STACK = docker compose -f docker-compose.yml
INTEGRATION_TEST_STACK = $(BASE_STACK) -f docker-compose-integration-test.yml
ALL_STACK = $(INTEGRATION_TEST_STACK)

# HELP =================================================================================================================
# This will output the help for each task
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker compose (without backend and reverse proxy)
	$(BASE_STACK) up --build -d db rabbitmq && docker compose logs -f
.PHONY: compose-up

compose-up-all: ### Run docker compose (with backend and reverse proxy)
	$(BASE_STACK) up --build -d
.PHONY: compose-up-all

compose-up-integration-test: ### Run docker compose with integration test
	$(INTEGRATION_TEST_STACK) up --build --abort-on-container-exit --exit-code-from integration-test
.PHONY: compose-up-integration-test

compose-down: ### Down docker compose
	$(ALL_STACK) down --remove-orphans
.PHONY: compose-down

swag-v1: ### swag init
	swag init -g internal/controller/http/router.go
.PHONY: swag-v1

deps: ### deps tidy + verify
	go mod tidy && go mod verify
.PHONY: deps

deps-audit: ### check dependencies vulnerabilities
	govulncheck ./...
.PHONY: deps-audit

format: ### Run code formatter
	gofumpt -l -w .
	gci write . --skip-generated -s standard -s default
.PHONY: format

run: deps swag-v1 ### swag run for API v1
	go mod download && \
	CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

test: ### run test
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

mock: ### run mockgen
	mockgen -source ./internal/repo/contracts.go -package usecase_test > ./internal/usecase/mocks_repo_test.go
	mockgen -source ./internal/usecase/contracts.go -package usecase_test > ./internal/usecase/mocks_usecase_test.go
.PHONY: mock

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations '$(word 2,$(MAKECMDGOALS))'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

bin-deps: ### install tools
	GOBIN=$(LOCAL_BIN) go install tool
.PHONY: bin-deps

pre-commit: swag-v1 mock format linter-golangci test ### run pre-commit
.PHONY: pre-commit

# Proto generation
PROTO_DIR=internal/controller/grpc/proto
PROTO_OUT=pkg/pb
PROTOC=protoc
PROTOC_OPTS=--go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative
GOBIN=$(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN=$(shell go env GOPATH)/bin
endif

.PHONY: proto-clean proto-all proto-install

proto-install: #### install protoc plugins
	@echo "Installing protoc plugins..."
	@echo "GOBIN: $(GOBIN)"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Verifying installations..."
	@test -f "$(GOBIN)/protoc-gen-go" || (echo "protoc-gen-go not found in $(GOBIN)" && exit 1)
	@test -f "$(GOBIN)/protoc-gen-go-grpc" || (echo "protoc-gen-go-grpc not found in $(GOBIN)" && exit 1)

proto-all: proto-install ### generate all proto files
	@echo "Generating all proto files..."
	@PATH="$(GOBIN):$$PATH" find $(PROTO_DIR) -name "*.proto" -exec sh -c ' \
		file_path="{}"; \
		rel_path=$${file_path#$(PROTO_DIR)/}; \
		service=$$(echo "$$rel_path" | cut -d/ -f1); \
		version=$$(echo "$$rel_path" | cut -d/ -f2); \
		file=$$(basename "$$file_path"); \
		echo "Processing $$service/$$version/$$file"; \
		mkdir -p $(PROTO_OUT)/$$service/$$version; \
		cd $(PROTO_DIR) && $(PROTOC) \
			--go_out=../../../../$(PROTO_OUT) \
			--go_opt=paths=source_relative \
			--go-grpc_out=../../../../$(PROTO_OUT) \
			--go-grpc_opt=paths=source_relative \
			$$rel_path \
	' \;

proto-clean: ### clean generated proto files
	@echo "Cleaning generated proto files..."
	@rm -rf $(PROTO_OUT)/*

proto: proto-install ### generate specific proto file
	@echo "Generating proto files..."
	@mkdir -p $(PROTO_OUT)/auth/v1
	@cd $(PROTO_DIR) && PATH="$(GOBIN):$$PATH" $(PROTOC) \
		--go_out=../../../../$(PROTO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=../../../../$(PROTO_OUT) \
		--go-grpc_opt=paths=source_relative \
		auth/v1/auth_service.proto
