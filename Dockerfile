# Use the official Go image as the base image
FROM golang:alpine3.19

# Set the working directory inside the container
WORKDIR /goapp

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .
COPY config_prod.yaml ./config.yaml

# Build the Go application
RUN go build -o main .

# EXPOSE 8080

# Set the entry point command to run the built binary
CMD ["./main"]