FROM golang:1.25-alpine as build

WORKDIR /src/dataproxyupdater
ADD . /src/dataproxyupdater
RUN CGO_ENABLED=0 go build -o /bin/dataproxyupdater ./cmd/dataproxyupdater

FROM ubuntu:24.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build /bin/dataproxyupdater /app/dataproxyupdater
EXPOSE 8000
CMD ["/app/dataproxyupdater"]
