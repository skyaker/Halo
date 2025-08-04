FROM golang:1.24.3-alpine AS builder

WORKDIR /usr/src/app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/category_service ./cmd/category_service

FROM alpine:latest

COPY --from=builder /usr/local/bin/category_service /usr/local/bin/category_service

CMD ["category_service"]
