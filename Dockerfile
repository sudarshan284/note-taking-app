# Start from a Go language base image
FROM golang:1.16-alpine

# Install Git
RUN apk update && apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Enable Go Modules and vendoring support
ENV GO111MODULE=on
ENV GOPROXY=direct
ENV GOSUMDB=off

# Build the Go application
RUN go mod download
RUN go build -o main .

# Expose port 8080 for the application
EXPOSE 8080

# Run the executable when the container starts
CMD ["./main"]
