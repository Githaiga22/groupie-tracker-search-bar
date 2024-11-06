# Step 1: Use the official Go image as a base image
FROM golang:1.22.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum if it exists, and download dependencies
COPY go.mod ./

RUN go mod download

# Copy the rest of the application code (source code)
COPY . .

# Build the Go application
RUN go build -o main .

# Step 2: Use Alpine as the base image for the final image
FROM alpine:latest

# Install necessary dependencies for running Go apps in Alpine
RUN apk --no-cache add ca-certificates libc6-compat

# Set the working directory to /app to match the application structure
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/main /app/main

# Copy the templates and static files from the builder stage into appropriate paths
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static

# Expose the port the app runs on
EXPOSE 8081

# Command to run the executable
CMD ["/app/main"]
