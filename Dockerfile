# To build the njuno image, just run:
# > docker build -t njuno .
#
# In order to work properly, this Docker container needs to have a volume that:
# - as source points to a directory which contains a config.yaml and firebase-config.yaml files
# - as destination it points to the /home folder
#
# Simple usage with a mounted data directory (considering ~/.njuno/config as the configuration folder):
# > docker run -it -v ~/.njuno/config:/home njuno njuno parse config.yaml firebase-config.json
#
# If you want to run this container as a daemon, you can do so by executing
# > docker run -td -v ~/.njuno/config:/home --name njuno njuno
#
# Once you have done so, you can enter the container shell by executing
# > docker exec -it njuno bash
#
# To exit the bash, just execute
# > exit
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev py-pip

# Set working directory for the build
WORKDIR /go/src/github.com/MonikaCat/njuno

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk update
RUN apk add --no-cache $PACKAGES && \
    make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /home

# Install bash
RUN apk add --no-cache bash

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/njuno /usr/bin/njuno