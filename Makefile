#!make

fmt:
	go fmt ./...

up:
	docker-compose up -d

logs:
	docker-compose logs -f go-stone