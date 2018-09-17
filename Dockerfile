FROM golang:1.10-alpine as base

RUN mkdir -p /go/src/github.com/janbaer/github-oauth-bridge
WORKDIR /go/src/github.com/janbaer/github-oauth-bridge
COPY . .
RUN pwd && ls -a
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main

FROM alpine:edge
RUN apk --update add ca-certificates
RUN mkdir /app
COPY --from=base /go/src/github.com/janbaer/github-oauth-bridge/main /app/github-oauth-bridge
COPY --from=base /go/src/github.com/janbaer/github-oauth-bridge/config.prod.json /app/
ENV ENV=PROD
EXPOSE 8080
WORKDIR /app
CMD ["/app/github-oauth-bridge"]
