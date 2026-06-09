# Multi-stage build for the Go consul-api proxy.
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
# Resolve dependencies (including transitive) and produce go.sum at build time.
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o consul-api .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/consul-api ./consul-api
EXPOSE 3000
CMD ["./consul-api"]
