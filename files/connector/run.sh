#!/bin/bash

CURRENT_DIR=`pwd`

main() {
    echo "Loading connector env"
    . $CURRENT_DIR/.env

    echo "Loading connector docker image"
    docker load < $CURRENT_DIR/roadrunner.docker

    mkdir "${CURRENT_DIR}/.rr_ssh"

    echo "Getting client certs from datahub"
    docker run -it -v "${CURRENT_DIR}/":"/tmp/config" \
           -v "${CURRENT_DIR}/.rr_ssh":"/go/bin/.ssh/roadrunner" gcr.io/npav-172917/adh-roadrunner:"${VERSION}" login --config=/tmp/config/adh-roadrunner.yml

    echo "Starting connector"
    docker run -d -v "${CURRENT_DIR}/":"/tmp/config" -v "${FILE_DIR}":"/tmp/files" -v "${CURRENT_DIR}/.rr_ssh":"/go/bin/.ssh/roadrunner" \
           --restart always --name aod-connector \
           --add-host "${HOST}:${IP}" \
           gcr.io/npav-172917/adh-roadrunner:"${VERSION}" data --config=/tmp/config/adh-roadrunner.yml

}

main "$@"
