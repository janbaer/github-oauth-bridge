test:
	@go test -cover ./...
.PHONY: test

install:
	@go get -d -u github.com/golang/dep
	@dep ensure
.PHONY: install

size:
	@curl -sL https://gist.githubusercontent.com/tj/04e0965e23da00ca33f101e5b2ed4ed4/raw/9aa16698b2bc606cf911219ea540972edef05c4b/gistfile1.txt | bash
.PHONY: size

deploy:
	@now && now alias
.PHONY: deploy

