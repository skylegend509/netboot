# Build Stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
# COPY go.sum ./ # No go.sum yet as no dependencies

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download 

COPY . .

# Build the application
RUN go build -o netboot ./cmd/netboot

# Run Stage
FROM alpine:latest  

WORKDIR /app

# Create directory for ISOs
RUN mkdir -p /app/isos /app/uploads

# Copy the binary from builder
COPY --from=builder /app/netboot .
# Copy static files
COPY --from=builder /app/static ./static

# Expose port
EXPOSE 8080

# Environment variables
ENV PORT=:8080
ENV ISO_DIR=/data/isos
ENV UPLOAD_DIR=/data/uploads

# Mount point for ISOs
VOLUME ["/data/isos"]

CMD ["./netboot"]
