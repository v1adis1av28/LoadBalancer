FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/load-balancer ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/load-balancer .
COPY pkg/models/config.json .

RUN chmod +x ./load-balancer

EXPOSE 8080

CMD ["./load-balancer"]