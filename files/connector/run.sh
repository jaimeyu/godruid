#!/bin/bash

CURRENT_DIR=`pwd`

main() {
    echo "Loading connector env"
    . $CURRENT_DIR/.env

    echo "Loading connector docker image"
    docker load < $CURRENT_DIR/roadrunner.docker

    echo "Starting connector"
    docker run --rm -d --name aod-connector -v "${CURRENT_DIR}/":"/tmp/config" -v "${FILE_DIR}":"/tmp/files" \
           --restart always \
                gcr.io/npav-172917/adh-roadrunner:"${FILE_DIR}" data --config=/tmp/config/connector-config.yml

}

main "$@"
