# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# Targeting cmd/main.go and naming the binary 'fleet-api'
RUN CGO_ENABLED=0 GOOS=linux go build -o fleet-api cmd/main.go

# Stage 2: Final minimal image
FROM alpine:latest

# Install CA certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/fleet-api .

# Copy environment file for local runtime (Note: in K8s, use ConfigMaps instead)
COPY --from=builder /app/.env* .

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["./fleet-api"]
