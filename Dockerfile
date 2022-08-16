FROM golang:1.18.3-alpine3.16 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o /golang-api

CMD ["/golang-api"]
