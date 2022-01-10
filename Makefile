# Makefile for build

APP=wgc

PLATFORMS=darwin linux windows
ARCHITECTURES=amd64 arm64

LDFLAGS=-ldflags '-s -w -extldflags "-static"' 


all: clean build_all

build:
	go build ${LDFLAGS} -o dist/${APP} bin/main.go

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build $(LDFLAGS) -o dist/$(APP)-$(GOOS)-$(GOARCH) bin/main.go)))
clean:
	@rm -rf dist

.PHONY: all build build_all clean
