test: export SECRET=Dvb6phNfnHUAjQAgGMm6KJMUXmZau3cS
test:
	@go test -v -cover ./...
.PHONY: test

watch-test:
	watchexec --exts go make test
.PHONY: watch-test

install:
	@npm install -g now
	@go mod download
.PHONY: install

run-dev:
	@go run main.go
.PHONY: run-dev

deploy:
	@now deploy --prod --local-config=now.prod.json --no-clipboard
.PHONY: deploy

