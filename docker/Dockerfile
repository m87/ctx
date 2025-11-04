# syntax=docker/dockerfile:1

# Build stage
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(if [ "$TARGETPLATFORM" = "linux/arm64" ]; then echo arm64; else echo amd64; fi) go build -o ctx-go ./main.go

# Final image
FROM --platform=$TARGETPLATFORM alpine:3.19
WORKDIR /app
COPY --from=builder /app/ctx-go ./ctx-go
EXPOSE 8080
ENTRYPOINT ["/app/ctx-go"]
