# Start from the latest golang base image
FROM golang:latest AS build

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Final stage
FROM alpine:latest  

# Set the working directory
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=build /app/main .

# Expose port 8888 to the outside world
EXPOSE 8888

# Command to run the executable
CMD ["./main"]