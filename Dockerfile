# Stage 1: Build (Native Architecture)
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

# Target architecture passed by buildx
ARG TARGETARCH
ARG ENTRY=monolith

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build natively for target architecture
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /server ./cmd/${ENTRY}

# Stage 2: Lean Runner
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Jakarta

WORKDIR /

COPY --from=builder /server /server
COPY .env .env

EXPOSE 8080

CMD ["/server"]
