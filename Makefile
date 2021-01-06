PWD = $(shell pwd -L)
IMAGE_NAME = gsabadini/go-bank-transfer
DOCKER_RUN = docker run --rm -it -w /app -v ${PWD}:/app -v ${GOPATH}/pkg/mod/cache:/go/pkg/mod/cache golang:1.14-stretch

init:
	cp .env.example .env

test:
	${DOCKER_RUN} go test -cover ./...

test-local:
	go test -cover ./...

coverage-report-func:
	${DOCKER_RUN} go test -covermode=set -coverprofile=coverage.out -p=1 ./... && \
 	go tool cover -func=coverage.out

coverage-report-browser:
	${DOCKER_RUN} go test -coverprofile coverage.out -p=1 ./... && \
	go tool cover -html=coverage.out -o coverage.html && \
	xdg-open ./coverage.html

coverage-report-text:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

code-review: fmt vet test

fmt:
	go fmt ./...

vet:
	go vet ./...

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f app

enter-container:
	docker-compose exec app bash

build-image:
	docker build -t ${IMAGE_NAME} -f Dockerfile .

run-image:
	docker run --rm -it --network -transfer_bank --env-file .env -p 3001:3001 ${IMAGE_NAME}

build:
	go build -a --installsuffix cgo --ldflags="-s" -o main

ci:
	docker run --rm -it -v $(PWD):/app -w /app golangci/golangci-lint:v1.24.0 \
	golangci-lint run
	--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...

.PHONY:ci coverage-report enter-container logs down up vet fmt init code-review test test-local build build-image