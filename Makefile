test:
	@go test -v -cover ./...
.PHONY: test

install:
	@npm install -g now
	@go mod download
.PHONY: install

run:
	@go run main.go
.PHONY: run

build:
	@docker build -t janbaer/oauth-bridge:latest .
.PHONY: build

deploy:
	@now --target production -e ENV=PROD
.PHONY: deploy

