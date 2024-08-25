# Stage1: BUILD
FROM golang:1.22 AS builder

LABEL maintainer="RespondNow"

ARG TARGETOS=linux
ARG TARGETARCH

# Set up the build environment
WORKDIR /respond

# Copy go.mod and go.sum files first, to leverage Docker layer caching for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the remaining project files
COPY . .

# Set environment variables for Go cross-compilation
ENV GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

RUN go env

# Install swag for generating swagger documentation
RUN GO111MODULE=on GOBIN=/usr/local/bin go install github.com/swaggo/swag/cmd/swag@v1.16.3

# Build the Go binary
RUN CGO_ENABLED=0 go build -o /output/respond -v ./cmd/

# Stage2: PACKAGING
# Use RedHat UBI minimal image as base
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4

LABEL maintainer="RespondNow"

# Set application directory
ENV APP_DIR="/respond"

# Copy the built binary from the builder stage
COPY --from=builder /output/respond $APP_DIR/
RUN chown 65534:0 $APP_DIR/respond && chmod 755 $APP_DIR/respond

# Set working directory and user permissions
WORKDIR $APP_DIR
USER 65534

# Run the application
CMD ["./respond"]

# Expose the application's port
EXPOSE 8080
