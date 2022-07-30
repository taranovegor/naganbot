FROM golang:1.20

WORKDIR /app

COPY ./ /app

RUN go mod download

ENTRYPOINT go run main.go
