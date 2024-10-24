# Use the Go image with Alpine for building and running the application
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the code and build the application
COPY . .
RUN go build -o app ./main.go

# Expose the port on which the app will run
EXPOSE 8081

# Command to run the binary
CMD ["./app"]