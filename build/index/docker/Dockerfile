# --- Build Stage ---
FROM golang:1.21-alpine as build

# Set the working directory inside the container for the build stage
WORKDIR /src/index

# Copy the entire project to the working directory
ADD . /src/index

# Build the Go app with CGO disabled to create a fully static binary,
# output the executable to /bin/index, compile the index app under ./cmd/index
RUN CGO_ENABLED=0 go build -o /bin/index ./cmd/index

# --- Runtime Stage ---
FROM ubuntu:20.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Copy the static binary from the build stage to the runtime stage
COPY --from=build /bin/index /app/index

EXPOSE 8000

CMD ["/app/index"]
