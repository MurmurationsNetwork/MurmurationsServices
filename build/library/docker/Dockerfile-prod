# --- Build Stage ---
FROM golang:1.19.5-alpine as build

# Set the working directory inside the container for the build stage
WORKDIR /src/library

# Copy the entire project to the working directory
ADD . /src/library

# Build the Go app with CGO disabled to create a fully static binary,
# output the executable to /bin/library, compile the library app under ./cmd/library
RUN CGO_ENABLED=0 go build -o /bin/library ./cmd/library

# --- Runtime Stage ---
FROM ubuntu:20.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Copy the static binary from the build stage to the runtime stage
COPY --from=build /bin/library /app/library

# Copy Static Files
COPY --from=build /src/library/services/library/static /app/static

EXPOSE 8000

CMD ["/app/library"]
