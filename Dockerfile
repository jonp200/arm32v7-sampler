# ---- Build Stage ----
FROM golang:1.20-bullseye AS builder

# Install ARM cross-compilation toolchain
RUN dpkg --add-architecture armhf && \
    apt-get update && \
    apt-get install -y \
        gcc-arm-linux-gnueabihf \
        libc6-dev-armhf-cross \
        ca-certificates

WORKDIR /app

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Set environment variables for Go cross-compilation
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CC=arm-linux-gnueabihf-gcc

# Enable CGO and statically link libraries
ENV CGO_ENABLED=1
ENV CGO_LDFLAGS='-static'

# Cross-compile for ARMv7 (32-bit)
RUN go build -o app .

# ---- Final Stage ----
FROM arm32v7/debian:bullseye-slim

WORKDIR /app

# Ensure CA certs are available in target image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy built binary
COPY --from=builder /app/app .

# Expose HTTP port
EXPOSE 8080

# Run the app
ENTRYPOINT ["./app"]
