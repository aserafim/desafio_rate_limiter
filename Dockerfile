# Dockerfile
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

# Se precisar de variáveis de ambiente
COPY .env .env

EXPOSE 8080

CMD ["./app"]
