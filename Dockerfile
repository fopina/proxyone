FROM alpine:3.20

COPY proxyone /usr/bin/proxyone
ENTRYPOINT ["/usr/bin/proxyone"]
