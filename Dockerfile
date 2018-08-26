FROM golang:1.10-alpine as base

RUN mkdir /go/src/oauthbridge
WORKDIR /go/src/oauthbridge
COPY . .
RUN ls -a
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main

FROM alpine:edge
RUN apk --update add ca-certificates
RUN mkdir /app && mkdir /app/config
COPY --from=base /go/src/oauthbridge/main /app/oauth-bridge
COPY --from=base /go/src/oauthbridge/config.prod.json /app/
ENV ENV=PROD
EXPOSE 8080
WORKDIR /app
CMD ["/app/oauth-bridge"]
