# services/location-service/Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum ./
COPY pkg/ ./pkg/
COPY gateway/ ./gateway/
COPY services/location-service/ ./services/location-service/

WORKDIR /app
RUN go build -o location-service ./services/location-service

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/location-service .

EXPOSE 8080

CMD ["./location-service"]