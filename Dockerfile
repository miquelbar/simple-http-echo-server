ARG GOLANG_VERSION=1.12
FROM golang:$GOLANG_VERSION as builder

RUN apt update && \
    apt install -y make git && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/simple-http-echo-server
ADD . .

RUN \
    make build-osx && \
    make build-linux

FROM scratch
COPY --from=builder /go/src/simple-http-echo-server/bin/simple-http-echo-server* /
ENTRYPOINT ["/simple-http-echo-server-linux"]
