#!/bin/sh

# Generates the API Server from a swagger file
#=============================================================

# Build up the full swagger file from the component files:
cd files
cat swagger/header.yml swagger/paths/*.yml swagger/definitions/header.yml swagger/definitions/definitions-*.yml > swagger.yml

cd ../
swagger generate server -f files/swagger.yml --exclude-main --exclude-spec -m swagmodels -A gather