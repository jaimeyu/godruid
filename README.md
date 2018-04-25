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

	1. `DOCKER_TAG`: The image tag to use (for instance: export `TAG=1.0.2; make docker` would build `gcr.io/npav-172917/adh-gather:1.0.2`)
	2. `DOCKER_IMAGE_NAME: Changes the name of the image (defaults to adh-gather)
	3. `DOCKER_REPO_NAME: Changes the repository to push the image to (defaults to gcr.io/npav-172917/)

You can generate the necessary gRPC Service, REST Reverse Proxy, and REST Swagger Definition by executing.

	1. ./generateFromProto.sh  


## Configuration

You can modify the following values in  the `config/adh-gath.yml` file:

	1. datastore
    	1. ip: the host of the CouchDB datastore. (Use full http:// format) 
    	2. port: the port used to communicate with CouchDB.
  	2. rest
    	1. ip: the host of the ADH-Gather REST API.
    	2. port: the port for the ADH-Gather REST API.
  	3. grpc
		1. ip: the host of the ADH-Gather gRPC Service.
    	2. port: the port for the ADH-Gather gRPC Service.
  	4. args
    	1. admindb: 
			- name: the name of the admin database used for adh-gather
			- impl: type of datastore to use for Admin Service. (0=InMemory, 1=CouchDB)
		2. tenantdb: type of datastore to use for Tenant Service. (0=InMemory, 1=CouchDB)
		3. pouchplugindb: type of datastore to use for PouchDB Plugin Service. (0=InMemory, 1=CouchDB)
		4. testdatadb: type of datastore to use for TestData Service. (0=InMemory, 1=CouchDB)
	5. kafka
		1. broker: the hostname and port of the kafka broker to listen to for change notifications

The following command line arguments are also available when executing the adh-gather program:
	
	1. config: Specify a configuration file to use (--config=path_to_config_file, default: "config/adh-gather.yml")
	2. tls: Specify if TLS should be enabled (--tls=true/false, default: true)
	3. tlskey: Specify a TLS Key file. (--tlskey=path_to_tls_key, default: "/run/secrets/tls_key")
	4. tlscert: Specify a TLS Cert file. (--tlscert=path_to_tls_cert, default: "/run/secrets/tls_crt")
	5. changeNotifications: Specify if Change Notifications should be enabled. (--changeNotifications=true/false, default: true)
