# Multi-stage build for the Go logging-api.
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -o /logging-api .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /logging-api /logging-api
EXPOSE 3000
ENTRYPOINT ["/logging-api"]
