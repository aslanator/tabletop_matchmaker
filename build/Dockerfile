# Start from golang base image
FROM golang:alpine as builder


# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
# Set the current working directory inside the container 
WORKDIR /app

ADD .. .
RUN go install github.com/ysmood/kit/cmd/guard@v0.25.11

RUN cat ./go.mod
# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Build the Go app
# RUN go build -a -installsuffix cgo -o main .

#Command to run the executable
ENTRYPOINT ["sh", "-c", "/go/bin/guard -w '**/*.go' -- go run ./bot.go"]
