FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
# Resolve and lock dependencies, then build a static binary.
RUN go mod tidy && CGO_ENABLED=0 go build -o tracing-api .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/tracing-api ./tracing-api
EXPOSE 3000
CMD ["./tracing-api"]
