#!make

code-review: fmt vet test

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -cover ./...

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f go-bank-transfer