FROM golang:1.18-alpine as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire source code to the working directory
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-todo cmd/todo/main.go

# Start a new stage from the Alpine Linux image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary file from the previous stage
COPY --from=builder /app/go-todo .
COPY --from=builder /app/.env .

# Run the application
CMD ["./go-todo"]
