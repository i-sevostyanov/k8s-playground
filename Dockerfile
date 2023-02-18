FROM golang:1.20.1-alpine3.17 AS build-env

# Get application name from args
ARG APP

# Set environment variables
ENV GO_WORKDIR /go/src/github.com/i-sevostyanov/k8s-playground
ENV GO111MODULE=on
ENV CGO_ENABLED=0

# Set working directory
WORKDIR $GO_WORKDIR
COPY . $GO_WORKDIR

# Add git and openssh
RUN set -eux; apk update; apk add --no-cache git openssh

# Install dependencies
RUN go mod download && go mod verify

# Build application
RUN go build -o /build/eny -ldflags "-s -w" cmd/$APP/*

# Put everything together in a clean image
FROM alpine:3.17

# Copy binary into PATH
COPY --from=build-env /build/eny /usr/local/bin/evy

ENTRYPOINT ["evy"]
