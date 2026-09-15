FROM golang:1.26.3-alpine AS builder

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o worker main.go

FROM alpine:3.23.3

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S worker-user && \
    adduser -S worker-user -G worker-user

WORKDIR /app

COPY --from=builder /api/worker /app/worker

USER worker-user

EXPOSE 8080

CMD ["./worker", "api"]