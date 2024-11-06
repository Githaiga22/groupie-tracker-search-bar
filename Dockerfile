# Step 1: Specify the base image
FROM golang:1.22.2 AS builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 4: Copy the entire project into the container's working directory
COPY . .

# Step 5: Build the application
RUN go build -o groupie-tracker-search-bar .

# Step 6: Use a lightweight base image for running the application
FROM gcr.io/distroless/base-debian11

# Step 7: Set the working directory and copy the built binary from the builder stage
WORKDIR /app
COPY --from=builder /app/groupie-tracker-search-bar .

# Step 8: Expose the port your application will run on
EXPOSE 8081

# Step 9: Set the entrypoint to run the application
CMD ["./groupie-tracker-search-bar"]
