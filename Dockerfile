FROM alpine:latest

RUN mkdir -p /go/bin

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

WORKDIR /go/bin

RUN apk add --no-cache ca-certificates 
RUN apk add --no-cache --virtual .build-deps librdkafka


COPY bin/alpine-adh-gather /go/bin/adh-gather

COPY config/ /config/
COPY files/ /files/

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
