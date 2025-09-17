# --- Build Stage ---
FROM golang:1.24-alpine as build

# Set the working directory inside the container for the build stage
WORKDIR /src/validation

# Copy the entire project to the working directory
ADD . /src/validation

# Build the Go app with CGO disabled to create a fully static binary,
# output the executable to /bin/validation, compile the validation app under ./cmd/validation
RUN CGO_ENABLED=0 go build -o /bin/validation ./cmd/validation

# --- Runtime Stage ---
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Copy the static binary from the build stage to the runtime stage
COPY --from=build /bin/validation /app/validation

EXPOSE 8000

CMD ["/app/validation"]
