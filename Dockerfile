# Use the official Golang base image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application files to the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd/main.go

# Expose the port your application will run on
EXPOSE 3000

# Command to run your application
CMD ["./main"]
