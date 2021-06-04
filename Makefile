ENVIRONMENT=development
SYSTEM=go-bank-transfer
SYSTEM_VERSION=$(shell git branch --show-current | cut -d '/' -f2)
PWD=$(shell pwd -L)
DOCKER_RUN=docker run --rm -it -w /app -v ${PWD}:/app -v ${GOPATH}/pkg/mod/cache:/go/pkg/mod/cache golang:1.16-buster

.PHONY: all
all: help
help: ## Display help screen
	@echo "Usage:"
	@echo "	make [COMMAND]"
	@echo "	make help \n"
	@echo "Commands: \n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
init: ## Create environment variables
	cp .env.example .env

.PHONY: test-local
test-local: ## Run local golang tests
	go test -cover -race ./...

.PHONY: test
test: ## Run golang tests
	${DOCKER_RUN} go test -cover -race ./...

.PHONY: test-report
test-report: ## Run tests with HTML coverage report
	${DOCKER_RUN} go test -covermode=count -coverprofile coverage.out -p=1 ./... && \
	go tool cover -html=coverage.out -o coverage.html && \
	xdg-open ./coverage.html

.PHONY: test-report-func
test-report-func: ## Run tests with func report -covermode=set
	${DOCKER_RUN} go test -covermode=set -coverprofile=coverage.out -p=1 ./... && \
	go tool cover -func=coverage.out

.PHONY: test-report-text
test-report-text:
	go test ./... -coverprofile=coverage.txt -covermode=atomic

# https://golangci-lint.run/usage/linters/
.PHONY: lint
lint: ## Lint with golangci-lint
	docker run --rm -it -v $(PWD):/app -w /app golangci/golangci-lint:v1.39-alpine \
	golangci-lint run \
	--exclude-use-default=false \
	--enable=gocyclo \
	--enable=bodyclose \
	--enable=goconst \
	--enable=sqlclosecheck \
	--enable=rowserrcheck \
	--enable=prealloc

.PHONY: fmt
fmt: ## Run go fmt
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: up
up: ## Run docker-compose up for creating and starting containers
	docker-compose up -d

.PHONY: down
down: ## Run docker-compose down for stopping and removing containers, networks, images, and volumes
	docker-compose down --remove-orphans

.PHONY: logs
logs: ## View container log
	docker-compose logs -f app

.PHONY: clean
clean: ## Clean build bin/
	@rm -rf bin/

.PHONY: build
build: clean ## Build golang project
	go build -o bin/$(SYSTEM) main.go

.PHONY: run
run: ## Run golang project
	go run main.go

.PHONY: docker-clean
docker-clean: ## Clean docker removes image
	docker rmi gsabadini/$(SYSTEM):$(SYSTEM_VERSION)

.PHONY: docker-build
docker-build: ## Build docker image for the project
	@docker build --target production -t gsabadini/$(SYSTEM):$(SYSTEM_VERSION) .

.PHONY: docker-run
docker-run: docker-deps ## Run docker container for image project
	docker run --rm -it \
	-e ENVIRONMENT=$(ENVIRONMENT) \
	-e SYSTEM=$(SYSTEM) \
	-e SYSTEM_VERSION=$(SYSTEM_VERSION) \
	-p 3001:3001 \
	--env-file .env \
	--network go-bank-transfer_bank  \
	--name $(SYSTEM)-$(SYSTEM_VERSION) gsabadini/$(SYSTEM):$(SYSTEM_VERSION)

docker-deps:
	docker-compose up -d postgres mongodb-primary mongodb-secondary mongodb-arbiter
	sleep 3