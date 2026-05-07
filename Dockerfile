FROM golang:1.25-alpine AS builder
ARG ENTRY=monolith
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/${ENTRY}

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Jakarta
COPY --from=builder /server /server
RUN touch .env
EXPOSE 8080
CMD ["/server"]
