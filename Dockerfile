# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git (needed for go modules).
RUN apk add --no-cache git

# Copy module files and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy source code.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kivo ./cmd/kivo

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder.
COPY --from=builder /app/kivo .

# Expose the Redis port.
EXPOSE 6379

# Run kivo.
CMD ["./kivo"]
