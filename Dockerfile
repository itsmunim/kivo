# syntax=docker/dockerfile:1

# --- Stage 1: build the web console (React SPA) ---
FROM node:22-alpine AS webui
WORKDIR /app/webui
COPY webui/package.json webui/package-lock.json* ./
RUN npm ci
COPY webui/ .
# Vite outputs to ../internal/webui/dist (i.e. /app/internal/webui/dist), which
# is where the Go builder stage copies the embedded UI from.
RUN npm run build
# Vite outputs to ../internal/webui/dist (embedded by Go).
RUN npm run build

# --- Stage 2: build the kivo binary with the UI embedded ---
FROM golang:1.26-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=webui /app/internal/webui/dist ./internal/webui/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o kivo ./cmd/kivo

# --- Stage 3: minimal runtime image ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/kivo .
EXPOSE 6379 3001
CMD ["./kivo"]