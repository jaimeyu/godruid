DOCKER_REPO_NAME:= gcr.io/npav-172917/
DOCKER_IMAGE_NAME := adh-gather
GO_REPOSITORY_PATH := github.com/accedian/$(DOCKER_IMAGE_NAME)
DOCKER_VER := $(if $(DOCKER_VER),$(DOCKER_VER),$(shell whoami)-dev)  
BIN_NAME := bin/alpine-$(DOCKER_IMAGE_NAME)
GO_SDK_IMAGE := gcr.io/npav-172917/docker-go-sdk
GO_SDK_VERSION := 1.0.1-alpine   
GOPATH := $(GOPATH)

SWAGGER_PATH := $(PWD)/files/swagger
SWAGGER_TMP_FILE := __swagger.yml
SWAGGER_TEMP := $(SWAGGER_PATH)/$(SWAGGER_TMP_FILE)
SWAGGER_FILES := $(SWAGGER_PATH)/header.yml \
    $(SWAGGER_PATH)/paths/paths-*.yml \
    $(SWAGGER_PATH)/definitions/header.yml \
    $(SWAGGER_PATH)/definitions/definitions-*.yml

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


swagger: $(SWAGGER_FILES)
	echo "Generating code from swagger files"
	rm -f $(SWAGGER_TEMP)
	cat $^ > $(SWAGGER_TEMP)
	docker run --rm -it -e GOPATH=$(GOPATH):/go -v$(HOME):$(HOME) -w $(PWD) quay.io/goswagger/swagger:0.15.0 generate server \
		 -f $(SWAGGER_TEMP) \
		 --exclude-main \
		 --exclude-spec \
		 -m swagmodels \
		 -A gather
	mv $(SWAGGER_TEMP) $(PWD)/files/swagger.yml
        

circleci-binaries:
	go build -o $(BIN_NAME) .

circleci-push: circleci-docker
	docker push $(DOCKER_REPO_NAME)$(DOCKER_IMAGE_NAME):$(DOCKER_VER)

circleci-docker: circleci-binaries
	docker build -t $(DOCKER_REPO_NAME)$(DOCKER_IMAGE_NAME):$(DOCKER_VER) .
	
.FORCE: 
clean:  
	rm -rf bin

