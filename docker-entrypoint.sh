#!/bin/sh
set -a  
 : ${SERVERMODE='true'} 
 : ${SERVER_DATASTORE_IP='http://couchdb'}
 : ${SERVER_CORS_ALLOWEDORIGINS='https://ui.*.npav.accedian.net'}


exec /go/bin/adh-gather --config /config/adh-gather.yml $@