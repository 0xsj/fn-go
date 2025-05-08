# gateway/Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum ./
COPY pkg/ ./pkg/
COPY gateway/ ./gateway/
COPY services/auth-service/ ./services/auth-service/
COPY services/entity-service/ ./services/entity-service/
COPY services/incident-service/ ./services/incident-service/
COPY services/location-service/ ./services/location-service/
COPY services/monitoring-service/ ./services/monitoring-service/
COPY services/notification-service/ ./services/notification-service/
COPY services/user-service/ ./services/user-service/

WORKDIR /app
RUN go build -o gateway ./gateway

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/gateway .

EXPOSE 8080

CMD ["./gateway"]