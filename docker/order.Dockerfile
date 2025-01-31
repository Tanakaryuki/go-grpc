FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order cmd/order/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/order .

EXPOSE 50053

CMD ["./order"]
