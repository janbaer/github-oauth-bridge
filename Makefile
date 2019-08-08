test:
	@go test -v -cover ./...
.PHONY: test

install:
	@go mod download
.PHONY: install

run:
	@go run main.go
.PHONY: run

build:
	@docker build -t janbaer/oauth-bridge:latest .
.PHONY: build

deploy:
	@now && now alias
.PHONY: deploy

