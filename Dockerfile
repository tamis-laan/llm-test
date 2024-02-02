# Use the official Golang base image for compilation
FROM golang:alpine as builder

# Install entr
RUN apk add fd entr make curl

# Set the working directory inside the container
WORKDIR /app

# Set the go binary directory
ENV GOBIN /usr/local/bin/

# Copy the Go module files to the container
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the binaries
RUN CGO_ENABLED=0 GOOS=linux go install ./cmd/...

###

# Use a minimal base image for the final container
FROM scratch as deployment

# Set the working directory inside the container
WORKDIR /app

# Copy executables from the builder stage
COPY --from=builder /usr/local/bin/* /usr/local/bin/

# Expose the port on which the application listens
EXPOSE 8080

# Start the Go Fiber application
CMD ["api"]

###

# Use a minimal base image for the final container
FROM alpine as ctl

# Set the working directory inside the container
WORKDIR /app

# Copy executables from the builder stage
COPY --from=builder /usr/local/bin/* /usr/local/bin/
