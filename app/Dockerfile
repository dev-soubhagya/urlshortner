# Use the official Golang image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY .env ./.env


# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o urlshortner .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./urlshortner"]
