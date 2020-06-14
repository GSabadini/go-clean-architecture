build:
	go build -a --installsuffix cgo --ldflags="-s" -o main

test:
	go test -cover ./...

test-container:
	docker-compose exec go-bank-transfer go test -cover ./...

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
	docker-compose logs -f go-bank-transfer

enter-container:
	docker-compose exec go-bank-transfer bash

lint-prepare:
	@echo "Installing golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run \
		--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...

ci:
	docker run --rm -it -v $(PWD):/app -w /app golangci/golangci-lint:v1.24.0 \
	golangci-lint run
	--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...