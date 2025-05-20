# Build stage
FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o numero

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/numero .
COPY README.md .
EXPOSE 8084
CMD ["/bin/sh", "-c", "./numero"]