FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]
