FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o inventory cmd/inventory/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/inventory .

EXPOSE 50054

CMD ["./inventory"]
