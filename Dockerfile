FROM alpine:3.23
# TODO: update binary name
COPY golang-template /usr/bin/golang-template
ENTRYPOINT ["/usr/bin/golang-template"]