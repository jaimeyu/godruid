#!/bin/sh
set -e
 : ${SERVERMODE='true'}
 : ${KAFKA_BROKER=kafka}
 : ${KAFKA_TOPIC=asm-data}
 : ${SPEC_FILE=/tmp/specFile.json}


exec /go/bin/adh-fedex --config /config/adh-fedex.yml $@