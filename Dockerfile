# Stage 1: Build
FROM golang:latest AS builder

WORKDIR /app

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Set environment variables for Go cross-compilation
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7

RUN go build -o app .

# Stage 2: Create minimal ARM image
FROM arm32v7/debian:12

# Copy built binary
COPY --from=builder /app/app /usr/local/bin/app

# Expose HTTP port
EXPOSE 8080

# Run the app
ENTRYPOINT ["/usr/local/bin/app"]
