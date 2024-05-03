FROM alpine:latest
RUN apk --no-cache add bash curl jq
COPY build/cleanup/docker/cleanup.sh /cleanup.sh
RUN chmod +x /cleanup.sh
ENTRYPOINT ["/cleanup.sh"]
