#!/bin/sh
set -a  
 : ${SERVERMODE='true'} 
 : ${SERVER_DATASTORE_IP='http://couchdb'}


exec /go/bin/adh-gather --config /config/adh-gather.yml --ingDict /files/defaultIngestionDictionary.json --swag /files/swagger.yml $@