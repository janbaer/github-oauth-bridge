FROM golang:1.12-alpine as base

RUN apk --update add git

RUN mkdir -p /github-oauth-bridge
WORKDIR /github-oauth-bridge

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main

FROM alpine:edge
RUN apk --update add ca-certificates
RUN mkdir /app
COPY --from=base /github-oauth-bridge/main /app/github-oauth-bridge
COPY ./config.prod.json /app/
ENV ENV=PROD
EXPOSE 8080
WORKDIR /app
CMD ["/app/github-oauth-bridge"]
