.PHONY: run build test swag migrate-up migrate-down docker-up docker-down

run:
	go run -ldflags="-linkmode=external" main.go

build:
	go build -ldflags="-linkmode=external" -o bin/emarket-api main.go

test:
	go test -v ./...

swag:
	swag init

docker-up:
	docker compose -f docker/docker-compose.yaml up -d

docker-down:
	docker compose -f docker/docker-compose.yaml down
