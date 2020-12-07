FROM golang:1.14.1-alpine3.11 AS builder

WORKDIR /ProjectNDPDNS

COPY ./ ./

RUN apk update \
    && apk add build-base gcc abuild binutils binutils-doc gcc-doc \
    && go build -a -o dns-main

FROM lsiobase/alpine:3.11

WORKDIR /npddns

EXPOSE 53
EXPOSE 53/udp

COPY --from=builder /ProjectNDPDNS/dns-main ./

ENTRYPOINT ["./dns-main"]
