# Start from golang base image
FROM golang:alpine as builder

# ENV GO111MODULE=on

# Add Maintainer info
LABEL maintainer="Erwin Salas <erwinsalas42@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container 
WORKDIR /app-order

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app-order/main .
COPY --from=builder /app-order/.env . 

# Install grpcurl
RUN apk --no-cache add curl
RUN wget -qO- https://github.com/fullstorydev/grpcurl/releases/download/v1.8.0/grpcurl_1.8.0_linux_x86_64.tar.gz | tar xvz -C /tmp && \
    mv /tmp/grpcurl /usr/local/bin/ && \
    rm -rf /tmp/grpcurl
    
# Expose port 50053 to the outside world
EXPOSE 50053

#Command to run the executable
CMD ["./main"]
