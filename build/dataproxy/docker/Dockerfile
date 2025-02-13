FROM golang:1.22-alpine as build

WORKDIR /src/dataproxy
ADD . /src/dataproxy
RUN CGO_ENABLED=0 go build -o /bin/dataproxy ./cmd/dataproxy/dataproxy

FROM golang:1.22-alpine as seeder-build

WORKDIR /src/dataproxy
ADD . /src/dataproxy
RUN CGO_ENABLED=0 go build -o /bin/seeder ./cmd/dataproxy/seeder
RUN CGO_ENABLED=0 go build -o /bin/ukseeder ./cmd/dataproxy/ukseeder

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build /bin/dataproxy /app/dataproxy
COPY --from=seeder-build /bin/seeder /app/seeder
COPY --from=seeder-build /bin/ukseeder /app/ukseeder
EXPOSE 8000
CMD ["/app/dataproxy"]
