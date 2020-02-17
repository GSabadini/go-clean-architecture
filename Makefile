#!make

fmt:
	go fmt ./...

test:
	go test -cover ./...

up:
	docker-compose up -d

logs:
	docker-compose logs -f go-stone