# --- Build Stage ---
# Use a Go base image with development tools for building
FROM golang:1.22-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to leverage Docker's build cache
COPY go.mod .
COPY go.sum .

# Download Go modules (dependencies)
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# CGO_ENABLED=0: Disables CGO, making the binary statically linked (no external C dependencies)
# GOOS=linux GOARCH=amd64: Ensures it's compiled for Linux AMD64 architecture (standard for most containers)
# -a -tags netgo: Builds a statically linked binary with network support
# -o go-sample-app: Specifies the output executable name
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -o go-sample-app .

# --- Production/Runtime Stage ---
# Use a minimal, secure, and actively supported Alpine base image for the final container
# alpine:3.20 is a good choice as of current knowledge, replacing EOL 3.12
FROM alpine:3.20

# Add Maintainer Info
LABEL maintainer="Community Engineering Team <community-engg@harness.io.>"

# Set the working directory for the application
WORKDIR /app

# Copy the compiled Go binary from the 'builder' stage into the final image
COPY --from=builder /app/go-sample-app /usr/local/bin/go-sample-app

# Expose the port your Go application listens on
EXPOSE 8080

# Command to run the executable when the container starts
# Use the full path to the executable
ENTRYPOINT ["/usr/local/bin/go-sample-app"]
