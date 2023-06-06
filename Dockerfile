
# Base image
FROM golang:1.18-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the application
RUN go build -o main .

# Expose the port on which the application will run
EXPOSE 8000

# Run the application
CMD ["./main"]