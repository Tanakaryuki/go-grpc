FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user cmd/user/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/user .

EXPOSE 50051

CMD ["./user"]
