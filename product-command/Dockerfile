FROM golang:1.26.5-alpine AS builder

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o product-command main.go

FROM alpine:3.23.3

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S product-command-user && \
    adduser -S product-command-user -G product-command-user

WORKDIR /app

COPY --from=builder /api/product-command /app/product-command

USER product-command-user

EXPOSE 8081

CMD ["./product-command", "api"]