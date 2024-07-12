# Stage 1: Build stage
# Use the official Golang image with Debian Bullseye as the base image
FROM golang:1.22-bullseye as builder

# Set the working directory inside the container to /app
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod .
COPY go.sum .

# Download and cache Go module dependencies
RUN go mod download

# Copy the entire source code to the working directory
COPY . .

# Set the Go build cache directory
ENV GOCACHE=/root/.cache/go-build

# Use a build cache for Go to speed up the build process
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o app ./cmd/app

# Stage 2: Final stage
# Use the official Ubuntu 22.04 image as the base image for the final stage
FROM ubuntu:22.04

# Create the /app directory inside the container
RUN mkdir /app

# Set the working directory inside the container to /app
WORKDIR /app

# Copy the built binary from the build stage to the final stage's /app directory
COPY --from=builder /app/app .

# Set the entrypoint command to run the app binary
ENTRYPOINT ["./app"]
