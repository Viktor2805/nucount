FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o main ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Create non-root user
RUN adduser -D -u 1001 appuser && \
    chown -R appuser:appuser /app

USER 1001

# Match the port your app listens on
EXPOSE 3000

CMD ["/app/main"]
