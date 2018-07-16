#!/bin/sh

# Generates the API Server from a swagger file
#=============================================================

    
swagger generate server -f files/swagger.yml --exclude-main --exclude-spec -m swagmodels -A gather