#!/bin/sh

# Generates the gRPC Service, REST reverse proxy, and Swagger
# definition file for the REST Service
#=============================================================

# Generate the gRPC Service
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  gathergrpc/commonModels.proto \
  gathergrpc/adminModels.proto  \
  gathergrpc/tenantModels.proto \
  gathergrpc/metricModels.proto \
  gathergrpc/gather.proto 
  
        

# Generate the REST handler for the gRPC service
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true,request_context=true:. \
  gathergrpc/gather.proto       
  