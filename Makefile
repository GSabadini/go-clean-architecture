#!make

fmt:
	go fmt ./...

test:
	go test -cover ./...

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f go-stone