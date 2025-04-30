# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o pausarr .

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pausarr .

# Install timezone data if needed
RUN apk --no-cache add tzdata

CMD ["./pausarr"]