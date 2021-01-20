.PHONY: build build-alpine clean test help default

BIN_NAME:=$(notdir $(shell pwd))
IMAGE_NAME := flamefatex/${BIN_NAME}
REMOTE_DOCKER_URI := flamefatex/${BIN_NAME}

# git信息
BRANCH := $(shell git branch | grep \* | cut -d ' ' -f2)
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_COMMIT := ${GIT_COMMIT}${GIT_DIRTY}
VERSION := $(shell git describe --tags)
ifeq "${VERSION}" ""
	VERSION := "notag"
endif

# GOPATH
ifndef GOPATH
    GOPATH=$(shell go env GOPATH)
endif
ifeq "$(findstring ${GOPATH}, $(shell pwd))" "${GOPATH}"
    # Found
    IN_GOPATH=true
else
    # Not found
    IN_GOPATH=false
endif

default: build

help:
	@echo 'Management commands:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make build-alpine    Compile optimized for alpine linux.'
	@echo '    make build-linux     Compile optimized for linux.'
	@echo '    make package         Build final docker image with just the go binary inside'
	@echo '    make tag             Tag image created by package with latest, git commit and version'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make push            Push tagged images to registry'
	@echo '    make clean           Clean the directory tree.'
	@echo '    make update-cookiecutter           Update the base config. e.g. gitlab-ci.yml, Dockerfile, Makefile...'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION} ${GIT_COMMIT}"
	go build -ldflags "-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT}" \
		-o bin/${BIN_NAME}

build-alpine:
	@echo "building ${BIN_NAME} ${VERSION} ${GIT_COMMIT}"
	go build -ldflags "-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT}" \
		-o bin/${BIN_NAME}

build-linux:
	@echo "building ${BIN_NAME} ${VERSION} ${GIT_COMMIT}"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags "-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT}" \
		-o bin/${BIN_NAME}

package:
	@echo "building image ${BIN_NAME} ${VERSION} ${GIT_COMMIT}"
	go mod vendor
	# 加快编译
	docker build \
		--build-arg APP_NAME=${BIN_NAME} \
		--build-arg VERSION=${VERSION} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		-t ${IMAGE_NAME}:latest .
	rm -rf vendor

tag:
	@echo "Tagging image ${BIN_NAME} ${VERSION} ${GIT_COMMIT}"
	docker tag $(IMAGE_NAME):latest $(REMOTE_DOCKER_URI):latest
#	docker tag $(IMAGE_NAME):latest $(REMOTE_DOCKER_URI):${BRANCH}
#	docker tag $(IMAGE_NAME):latest $(REMOTE_DOCKER_URI):${VERSION}

push: tag
	docker push $(REMOTE_DOCKER_URI):latest
#	docker push $(REMOTE_DOCKER_URI):${BRANCH}
#	docker push $(REMOTE_DOCKER_URI):${VERSION}

local_run: package
	mkdir -p /tmp/$(BIN_NAME)
	docker run -p8084:8084 -v/tmp/$(BIN_NAME):/etc/${BIN_NAME} -it --rm $(IMAGE_NAME):latest

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}
	@test ! -e /tmp/${BIN_NAME} || rm -rf /tmp/${BIN_NAME}

update-cookiecutter:
	rm -rf /tmp/cookiecutter
	mkdir -p /tmp/cookiecutter
	cd /tmp/cookiecutter && cookiecutter https://github.com/flamefatex/cookiecutter-golang.git --no-input app_name=$(BIN_NAME)
	cp /tmp/cookiecutter/$(BIN_NAME)/Makefile . | true
	cp /tmp/cookiecutter/$(BIN_NAME)/Dockerfile . | true
	cp /tmp/cookiecutter/$(BIN_NAME)/.gitlab-ci.yml . | true
	rm -rf /tmp/cookiecutter

test:
	go test -v `go list ./... | grep -v vendor`

