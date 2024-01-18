# --- Build Stage ---
FROM golang:1.21-alpine as build

# Set the working directory inside the container for the build stage
WORKDIR /src/geoip

# Copy the entire project to the working directory
ADD . /src/geoip

# Build the Go app with CGO disabled to create a fully static binary,
# output the executable to /bin/geoip, compile the geoip app under ./cmd/geoip
RUN CGO_ENABLED=0 go build -o /bin/geoip ./cmd/geoip

# --- Runtime Stage ---
FROM ubuntu:20.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Copy the static binary from the build stage to the runtime stage
COPY --from=build /bin/geoip /app/geoip

#  Copy the GeoLite2-City.mmdb file
COPY --from=build /src/geoip/services/geoip/global/GeoLite2-City.mmdb /app/GeoLite2-City.mmdb

EXPOSE 8000

CMD ["/app/geoip"]
