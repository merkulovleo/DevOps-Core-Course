# Stage 1: Builder - compile the Go application
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy go mod file first for better layer caching
COPY go.mod ./

# Download dependencies (none in this case, but good practice)
RUN go mod download

# Copy source code
COPY main.go .

# Build the binary with optimizations
# CGO_ENABLED=0 creates a static binary
# -ldflags="-s -w" strips debug info for smaller binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server main.go


# Stage 2: Runtime - minimal scratch image
FROM scratch AS runtime

# Copy the binary from builder stage
COPY --from=builder /app/server /server

# Expose application port
EXPOSE 8080

# Run as the binary directly (scratch has no shell)
ENTRYPOINT ["/server"]
