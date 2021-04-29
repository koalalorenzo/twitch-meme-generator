BUILD_TARGET ?=

DOCKER_IMAGE ?= registry.gitlab.com/koalalorenzo/twitch-meme-generator
DOCKER_TAG ?= $(shell git log -1 --pretty=format:"%h")

ifeq (${BUILD_TARGET},rpi)
GOARCH := arm
GOOS := linux
GOARM=7
endif

.EXPORT_ALL_VARIABLES:

build:
	go build -o build/koalalorenzo-meme-generator 
	cp -R ./assets ./build/assets

run:
	go run ./*.go --webhook-enable
.PHONY: run

install:
	go install 

clean:
	rm -rf build
.PHONY: clean

#### Docker targets
docker_build:
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
.PHONY: docker_build

docker_push: docker_build
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
.PHONY: docker_push