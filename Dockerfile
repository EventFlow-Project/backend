FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -g '' appuser

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations

RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080

CMD ["./main"] 