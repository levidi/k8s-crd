# Use the official Golang image as a parent image
FROM golang:1.22.3 as builder

# Set the current working directory inside the container
WORKDIR /go/src/app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o manager .

# Use a minimal base image to package the final binary
FROM alpine:latest

# Set the working directory to /root/
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /go/src/app/manager .

# Command to run the executable
CMD ["./manager"]
