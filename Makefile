.PHONY: generate
generate:
	buf generate

.PHONY: server
server:
	go run cmd/main.go