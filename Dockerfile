FROM iron/go


ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

WORKDIR /go/bin

RUN  apk add --no-cache ca-certificates && \
     apk add --no-cache --virtual .build-deps librdkafka

RUN  rm -rf /go/pkg /go/src /usr/local/go

COPY bin/alpine-adh-gather /go/bin/adh-gather

COPY config/ /config/

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
