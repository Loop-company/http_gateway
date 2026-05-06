FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/http_gateway ./cmd

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /bin/http_gateway /app/http_gateway

EXPOSE 8080
CMD ["/app/http_gateway"]
