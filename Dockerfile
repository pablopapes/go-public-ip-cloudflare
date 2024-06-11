# Use the official Golang image from the Docker Hub
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD (Present Working Directory) in the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o main .

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]