# Use Golang with Alpine for a lightweight build
FROM golang:1.23.4-alpine3.19 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go application (fixing the path to main.go)
RUN go build -o main ./cmd/main.go

# Use a minimal image for the final container
FROM alpine:3.19

# Set working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file into the final image
COPY --from=builder /app/.env /

# Expose the port your Go app runs on
EXPOSE 8080

# Run the Go binary
CMD ["./main"]