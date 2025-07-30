FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth_notification ./cmd/main.go

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/auth_notification /auth_notification
COPY --from=builder /app/.env /.env

EXPOSE 8080

CMD ["/auth_notification"]