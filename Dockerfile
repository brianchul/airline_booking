# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build API Gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-gateway ./cmd/api-gateway

# Build Booking Service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o booking-service ./cmd/booking-service

# API Gateway stage
FROM alpine:latest AS api-gateway

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/api-gateway .

# Copy config files
COPY --from=builder /app/internal/config ./config

EXPOSE 8080

CMD ["./api-gateway"]

# Booking Service stage
FROM alpine:latest AS booking-service

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/booking-service .

# Copy config files
COPY --from=builder /app/internal/config ./config

EXPOSE 8081

CMD ["./booking-service"]