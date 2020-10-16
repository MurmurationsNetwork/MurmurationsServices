FROM golang:1.15-buster as build

ADD . /src/indexer
WORKDIR /src/indexer
RUN go build -o /bin/indexer ./indexer/cmd/indexer

FROM ubuntu:16.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=build /bin/indexer /app/
EXPOSE 8080
ENTRYPOINT ["/app/indexer"]
CMD ["serve"]
