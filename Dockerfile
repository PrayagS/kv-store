FROM golang:1.17.0-buster as builder
WORKDIR /kv-store
COPY . .
RUN go build -o main ./cmd/main.go

FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /kv-store/main .
CMD ["./main"]
