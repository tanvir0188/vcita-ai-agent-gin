# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO required for go-sqlite3 (GORM SQLite driver)
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o vcita-agent ./cmd/server

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Non-root user (principle of least privilege)
RUN useradd -r -s /bin/false vcita
WORKDIR /app
RUN mkdir -p data logs && chown -R vcita:vcita /app

COPY --from=builder /app/vcita-agent .

USER vcita
EXPOSE 8443
CMD ["./vcita-agent"]
