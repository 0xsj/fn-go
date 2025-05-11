FROM golang:1.24-alpine AS builder
WORKDIR /app

# Define build argument for the service name
ARG SERVICE_NAME=auth-service

# Copy go.work files first to leverage layer caching
COPY go.work go.work.sum ./

# Copy required packages
COPY pkg/ ./pkg/
COPY gateway/ ./gateway/
COPY services/ ./services/

# Build the service using the SERVICE_NAME arg
RUN go build -o service ./services/${SERVICE_NAME}/cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/service .

EXPOSE 8080
CMD ["./service"]