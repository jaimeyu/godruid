# adh-gather

The magic circle component, responsible for gathering data from various datasource and presenting it to UI and other 3rdParties

## Get IT

Execute the the following

```go
go get github.com/accedian/adh-gather
```

This will clone and build the code in `$GOPATH/src/github.com/accedian/adh-gather`

## Build

Execute the following makefile targets to build this docker image

	1. `make docker`: builds the docker image
	2. `make push` : builds the docker image and push it to the gcr repository

You can set the following ENV variable before calling make:

	1. `DOCKER_TAG`: The image tag to use (for instance: export `TAG=1.0.2; make docker` would build `gcr.io/npav-172917/druid:1.0.2`)
	2. `DOCKER_IMAGE_NAME: Changes the name of the image (defaults to druid)
	3. `DOCKER_REPO_NAME: Changes the repository to push the image to (defaults to gcr.io/npav-172917/

You can generate the necessary gRPC Service, REST Reverse Proxy, and REST Swagger Definition by executing

	1. ./generateFromProto.sh  


## Configuration

You can modify the following values in  the `config/adh-gath.yml` file:

	1. Datastore
    	1. bindIP: the host of the CouchDB datastore. (Use full http:// format) 
    	2. bindPort: the port used to communicate with CouchDB.
  	2. REST
    	1. bindIP: the host of the ADH-Gather REST API.
    	2. bindPort: the port for the ADH-Gather REST API.
  	3. GRPC
		1. bindIP: the host of the ADH-Gather gRPC Service.
    	2. bindPort: the port for the ADH-Gather gRPC Service.
  	4. StartupArgs
    	1. adminDB: 
			- name: the name of the admin database used for adh-gather
			- impl: type of datastore to use for Admin Service. (0=InMemory, 1=CouchDB)
		2. tenantDB: type of datastore to use for Tenant Service. (0=InMemory, 1=CouchDB)
		3. pouchPluginDB: type of datastore to use for PouchDB Plugin Service. (0=InMemory, 1=CouchDB)

