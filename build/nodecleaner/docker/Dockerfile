# --- Build Stage ---
FROM golang:1.21-alpine as build

# Set the working directory inside the container for the build stage
WORKDIR /src/nodecleaner

# Copy the entire project to the working directory
ADD . /src/nodecleaner

# Build the Go app with CGO disabled to create a fully static binary,
# output the executable to /bin/nodecleaner, compile the nodecleaner app under ./cmd/nodecleaner
RUN CGO_ENABLED=0 go build -o /bin/nodecleaner ./cmd/nodecleaner

# --- Runtime Stage ---
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Copy the static binary from the build stage to the runtime stage
COPY --from=build /bin/nodecleaner /app/nodecleaner

EXPOSE 8000

CMD ["/app/nodecleaner"]
