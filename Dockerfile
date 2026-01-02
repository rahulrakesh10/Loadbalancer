# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the load balancer
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/load-balancer ./main.go

# Build backend server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/backend-server ./cmd/backend/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binaries from builder
COPY --from=builder /app/load-balancer .
COPY --from=builder /app/backend-server .
COPY --from=builder /app/config/servers.json ./config/

EXPOSE 8080 9001 9002 9003

CMD ["./load-balancer"]


