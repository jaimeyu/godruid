#!/bin/sh

# Generates the API Server from a swagger file
#=============================================================

# Build up the full swagger file from the component files:
cd files
cat swagger/header.yml swagger/paths/*.yml swagger/definitions/header.yml swagger/definitions/definitions-*.yml > swagger.yml

cd ../

# Remove the old generated files in case some are no longer needed.
rm -rf swagmodels
rm -rf restapi/operations

# Generate the server and models
swagger generate server -f files/swagger.yml --exclude-main --exclude-spec -m swagmodels -A gather