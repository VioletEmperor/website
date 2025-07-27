# syntax=docker/dockerfile:1

# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o website

# Stage 2: Create the final lean image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/website .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
EXPOSE 80
CMD ["./website"]