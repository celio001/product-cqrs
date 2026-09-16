FROM golang:1.26.5-alpine AS builder

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o product-query main.go

FROM alpine:3.23.3

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S product-query-user && \
    adduser -S product-query-user -G product-query-user

WORKDIR /app

COPY --from=builder /api/product-query /app/product-query

USER product-query-user

EXPOSE 8082

CMD ["./product-query", "api"]