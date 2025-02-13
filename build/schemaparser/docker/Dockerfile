# --- Build Stage ---
FROM golang:1.22-alpine as build

# Set the working directory inside the container for the build stage
WORKDIR /src/schemaparser

# Copy the entire project to the working directory
ADD . /src/schemaparser

# Build the Go app with CGO disabled to create a fully static binary,
# output the executable to /bin/schemaparser, compile the schemaparser app under ./cmd/schemaparser
RUN CGO_ENABLED=0 go build -o /bin/schemaparser ./cmd/schemaparser

# --- Runtime Stage ---
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Copy the static binary from the build stage to the runtime stage
COPY --from=build /bin/schemaparser /app/schemaparser

EXPOSE 8000

CMD ["/app/schemaparser"]
