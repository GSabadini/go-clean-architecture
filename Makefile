#!make

dependencies:
	go mod download

code-review: fmt vet test

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint ./...

test:
	go test -cover ./...

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f go-bank-transfer

enter-container:
	docker-compose exec go-bank-transfer bash