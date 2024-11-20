FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/cmd/api

RUN go build -o api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/cmd/api/api .
COPY --from=builder /app/.env .

CMD ["./api"]