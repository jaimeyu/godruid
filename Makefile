DOCKER_REPO_NAME:= gcr.io/npav-172917/
DOCKER_IMAGE_NAME := adh-gather
GO_REPOSITORY_PATH := github.com/accedian/$(DOCKER_IMAGE_NAME)
DOCKER_VER := $(if $(DOCKER_VER),$(DOCKER_VER),latest)  
BIN_NAME := bin/alpine-$(DOCKER_IMAGE_NAME)
GO_SDK_IMAGE := gcr.io/npav-172917/docker-go-sdk
GO_SDK_VERSION := 1.0.1-alpine
  
GOPATH := $(GOPATH)
all: docker

dockerbin: .FORCE
	echo "PATH is $(GOPATH)"
	docker run -it --rm \
		-e GOPATH=/root/go \
		-v "$(GOPATH):/root/go" \
		-w "/root/go/src/$(GO_REPOSITORY_PATH)" \
		$(GO_SDK_IMAGE):$(GO_SDK_VERSION) go build -o $(BIN_NAME) . 

docker: dockerbin
	 docker build -t $(DOCKER_REPO_NAME)$(DOCKER_IMAGE_NAME):$(DOCKER_VER) .

push: docker
	docker push $(DOCKER_REPO_NAME)$(DOCKER_IMAGE_NAME):$(DOCKER_VER)

circleci-binaries:
	go build -o $(BIN_NAME) .

circleci-push: circleci-docker
	docker push $(DOCKER_REPO_NAME)$(DOCKER_IMAGE_NAME):$(DOCKER_VER)

circleci-docker: circleci-binaries
	docker build -t $(DOCKER_REPO_NAME)$(DOCKER_IMAGE_NAME):$(DOCKER_VER) .
	
.FORCE: 
clean:  
	rm -rf bin

