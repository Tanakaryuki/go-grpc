FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o product cmd/product/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/product .

EXPOSE 50052

CMD ["./product"]
