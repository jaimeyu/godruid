#!/bin/bash

CURRENT_DIR=`pwd`

main() {
    echo "Loading connector env"
    . $CURRENT_DIR/.env

    echo "Loading connector docker image"
    docker load < $CURRENT_DIR/roadrunner.docker

    mkdir -p "${CURRENT_DIR}/.rr_ssh"

    echo "Getting client certs from datahub"
    docker run -it -v "${CURRENT_DIR}/":"/tmp/config" \
           --add-host "${DEPLOYMENT_HOSTNAME}:${DEPLOYMENT_IP}" \
           --add-host "${TENANT_HOSTNAME}:${TENANT_IP}" \
           -v "${CURRENT_DIR}/.rr_ssh":"/go/bin/.ssh/roadrunner" gcr.io/npav-172917/adh-roadrunner:"${VERSION}" login --config=/tmp/config/adh-roadrunner.yml

    echo "Stopping old connector"
    docker rm -f aod-connector-for-${TENANT_HOSTNAME}

    echo "Starting connector"
    docker run -d -v "${CURRENT_DIR}/":"/tmp/config" -v "${FILE_DIR}":"/tmp/files" -v "${CURRENT_DIR}/.rr_ssh":"/go/bin/.ssh/roadrunner" \
           --restart always \
           --name aod-connector-for-${TENANT_HOSTNAME} \
           --add-host "${DEPLOYMENT_HOSTNAME}:${DEPLOYMENT_IP}" \
           --add-host "${TENANT_HOSTNAME}:${TENANT_IP}" \
           gcr.io/npav-172917/adh-roadrunner:"${VERSION}" data --config=/tmp/config/adh-roadrunner.yml

}

main "$@"
