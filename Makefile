test:
	@go test -v -cover ./...
.PHONY: test

install:
	@go get -d -u github.com/golang/dep
	@dep ensure
.PHONY: install

run:
	@go run main.go
.PHONY: run

deploy:
	@now && now alias
.PHONY: deploy

