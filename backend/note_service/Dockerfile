FROM golang:1.24.3-alpine AS builder

WORKDIR /usr/src/app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/note_service ./cmd/note_service

FROM alpine:latest

COPY --from=builder /usr/local/bin/note_service /usr/local/bin/note_service

CMD ["note_service"]
