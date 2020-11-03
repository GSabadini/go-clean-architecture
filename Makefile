build:
	go build -a --installsuffix cgo --ldflags="-s" -o main

test:
	go test -cover ./...

test-container:
	docker-compose exec app go test -cover ./...

dependencies:
	go mod download

code-review: fmt vet test

init:
	cp .env.example .env

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

coverage-report:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

ci:
	docker run --rm -it -v $(PWD):/app -w /app golangci/golangci-lint:v1.24.0 \
	golangci-lint run
	--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...