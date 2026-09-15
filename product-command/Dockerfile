FROM golang:1.26.3-alpine AS builder

WORKDIR /api

FROM alpine:3.23.3

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S product-command-user && \
    adduser -S product-command-user -G product-command-user

WORKDIR /app

COPY --from=builder /api/product-command /app/product-command

USER product-command-user

EXPOSE 8080

CMD ["./product-command", "api"]