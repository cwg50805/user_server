# Use the official Golang image
FROM golang:1.22-alpine AS build

# Set the working directory in the container
WORKDIR /app

# Copy the code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Install MySQL and Redis clients
RUN apk --no-cache add mysql-client redis

# Copy the built executable from the build stage
COPY --from=build /app/main /app/main

# Set the working directory in the container
WORKDIR /app

# Run the Go application
CMD ["./main"]
