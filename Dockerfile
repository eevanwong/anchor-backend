# Use the official Golang image as the build environment
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Copy the .env file into the container
COPY .env .env

# Build the Go app (CGO_ENABLED=0 disables C code access and GOOS=linux enables code compilation in linux system)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage from scratch
FROM alpine:latest  

# # Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the .env file to the final image (to ensure app can reaad .env during runtime)
COPY --from=builder /app/.env ./

# Command to run the executable
CMD ["./main"]
