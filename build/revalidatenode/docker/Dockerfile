# --- Build Stage ---
FROM golang:1.19.5-alpine as build

# Set the working directory inside the container for the build stage
WORKDIR /src/revalidatenode

# Copy the entire project to the working directory
ADD . /src/revalidatenode

# Build the Go app with CGO disabled to create a fully static binary,
# output the executable to /bin/revalidatenode, compile the revalidatenode app under ./cmd/revalidatenode
RUN CGO_ENABLED=0 go build -o /bin/revalidatenode ./cmd/revalidatenode

# --- Runtime Stage ---
FROM ubuntu:20.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Copy the static binary from the build stage to the runtime stage
COPY --from=build /bin/revalidatenode /app/revalidatenode

EXPOSE 8000

CMD ["/app/revalidatenode"]
