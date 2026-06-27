# Build stage
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o GeminiWatermarkTool-Go main.go

# Run stage
FROM alpine:latest
RUN apk add --no-cache ffmpeg ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/GeminiWatermarkTool-Go .
ENTRYPOINT ["./GeminiWatermarkTool-Go"]
