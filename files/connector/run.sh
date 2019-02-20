#!/bin/bash
CURRENT_DIR=`pwd`

main() {
    echo "Loading connector env"
    . $CURRENT_DIR/.env

    echo "Loading connector docker image"
    docker load < $CURRENT_DIR/roadrunner.docker

    mkdir -p "${CURRENT_DIR}/.rr_ssh"

    proxy_option=
    if [[ ${HTTPS_PROXY} ]]; then
        proxy_option="-e https_proxy=${HTTPS_PROXY}"
    fi

    echo "Getting client certs from datahub"
    docker run -it -v "${CURRENT_DIR}/":"/tmp/config":z \
           ${proxy_option} \
           --add-host "${DEPLOYMENT_HOSTNAME}:${DEPLOYMENT_IP}" \
           --add-host "${TENANT_HOSTNAME}:${TENANT_IP}" \
           -v "${CURRENT_DIR}/.rr_ssh":"/go/bin/.ssh/roadrunner":z gcr.io/npav-172917/adh-roadrunner:"${VERSION}" login --config=/tmp/config/adh-roadrunner.yml

    echo "Stopping old connector"
    docker rm -f aod-connector-for-${TENANT_HOSTNAME}

    echo "Starting connector"
    docker run -d -v "${CURRENT_DIR}/":"/tmp/config":z -v "${FILE_DIR}":"/tmp/files":z -v "${CURRENT_DIR}/.rr_ssh":"/go/bin/.ssh/roadrunner":z \
           -v /proc:/prochost:z \
           -v /etc:/etchost:z \
           -v /var:/varhost:z \
           -v /sys:/syshost:z \
           --restart always \
           --name aod-connector2-for-${TENANT_HOSTNAME} \
           ${proxy_option} \
           --add-host "${DEPLOYMENT_HOSTNAME}:${DEPLOYMENT_IP}" \
           --add-host "${TENANT_HOSTNAME}:${TENANT_IP}" \
           --env-file "$CURRENT_DIR/.env" \
           -e HOST_PROC=/prochost \
           -e HOST_SYS=/syshost \
           -e HOST_ETC=/etchost \
           -e HOST_VAR=/varhost \
           gcr.io/npav-172917/adh-roadrunner:"${VERSION}" data --config=/tmp/config/adh-roadrunner.yml
}

main "$@"
