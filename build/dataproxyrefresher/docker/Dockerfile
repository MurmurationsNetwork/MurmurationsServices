# Using the official golang image based on Alpine for the build stage.
FROM golang:1.19.5-alpine as build

WORKDIR /src/dataproxyrefresher

ADD . /src/dataproxyrefresher

# Build the Go application with CGO disabled and output the binary to /bin/dataproxyrefresher.
RUN CGO_ENABLED=0 go build -o /bin/dataproxyrefresher ./cmd/dataproxyrefresher

# Using Ubuntu 20.04 as the base image for the final container.
FROM ubuntu:20.04

# Update the package list and install the ca-certificates package for SSL/TLS support.
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Copying the compiled binary from the build container to the final container.
COPY --from=build /bin/dataproxyrefresher /app/dataproxyrefresher

# Exposing port 8000 for the service to be accessible.
EXPOSE 8000

# Setting the default command to run the compiled Go application.
CMD ["/app/dataproxyrefresher"]
