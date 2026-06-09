# Multi-stage build: compile a static binary, then ship it on a tiny base image.
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
# Resolve dependencies and generate go.sum reproducibly inside the build.
RUN go mod download && go mod tidy
COPY cmd ./cmd
COPY internal ./internal
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /metrics-api ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /metrics-api /metrics-api
EXPOSE 3000
ENTRYPOINT ["/metrics-api"]
