FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY pkg/ ./pkg/
COPY gateway/ ./gateway/
COPY services/ ./services/
WORKDIR /app
RUN go build -o entity-service ./services/entity-service
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/entity-service .
EXPOSE 8080
CMD ["./entity-service"]