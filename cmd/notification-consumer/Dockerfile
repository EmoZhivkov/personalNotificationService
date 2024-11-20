FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/cmd/notification-consumer

RUN go build -o notification-consumer

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/cmd/notification-consumer/notification-consumer .
COPY --from=builder /app/.env .

CMD ["./notification-consumer"]