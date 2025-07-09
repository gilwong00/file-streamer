.PHONY: generate
generate:
	buf generate

.PHONY: server
server:
	go run cmd/main.go

docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down -v --remove-orphans