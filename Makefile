.PHONY : run clean build

##############################################################################
#
#  Makefile to manage a go app lifecycle
#
##############################################################################

VERSION ?= $(shell cat VERSION)

default: run

install:
	dep ensure
	sh ./init.sh

run:
	GO_ENV=local go run suitesync.go

release: clean
	curl -sL https://git.io/goreleaser | bash

tag-version:
	git tag ${VERSION} && git push origin ${VERSION}

build: clean
	go build -o suitesync suitesync.go

clean:
	rm -rf suitesync dist

test:
	go test ./... -cover -coverprofile=c.out -tags test
	go tool cover -html=c.out -o coverage.html
