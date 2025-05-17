FROM golang:1.21 AS build

WORKDIR /app

COPY ./ /app

RUN go mod download

ARG VERSION

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.Version=${VERSION}'" -o bin

FROM debian:12-slim

RUN apt-get update && apt-get install -y \
    ca-certificates

WORKDIR /naganbot

COPY --from=build /app/bin bin

ENTRYPOINT ["/naganbot/bin"]
