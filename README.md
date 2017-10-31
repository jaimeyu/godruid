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

