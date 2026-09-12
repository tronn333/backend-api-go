# syntax=docker/dockerfile:1

# ── Build stage ──────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

WORKDIR /src

# Download module dependencies first for better layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source and build a statically-linked binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/api

# ── Runtime stage ────────────────────────────────────────────────────
FROM alpine:3.19

# Run the app as an unprivileged user.
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=builder /out/server ./server

USER app
EXPOSE 8080

# The binary is PID 1 and handles SIGTERM/SIGINT for graceful shutdown.
HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
  CMD wget -q -O /dev/null http://localhost:8080/health || exit 1

ENTRYPOINT ["./server"]