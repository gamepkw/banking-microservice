# Use a multi-stage build to keep the final image minimal
FROM golang:1.21

# Set the working directory to /app
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.* ./

# Download dependencies
RUN go mod download

# Copy the local package files to the container's working directory
COPY . .

# Build the binary
RUN go build -o main

# Set the working directory to /app
WORKDIR /app

# Expose the port on which the application will run (adjust as needed)
EXPOSE 8090

# Command to run the executable
CMD ["./main"]
