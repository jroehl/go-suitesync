.PHONY : run clean build

##############################################################################
#
#  Makefile to manage a go app lifecycle
#
##############################################################################

default: run

install:
	dep ensure
	sh ./init.sh

run:
	go run suitesync.go

release:
	curl -sL https://git.io/goreleaser | bash

build: clean
	go build -o suitesync suitesync.go

pack-restlet: clean
	tar -czf restlet.tar.gz ./restlet/project

clean:
	rm -rf suitesync restlet.tar.gz
# test-unit:
# 	go test ./...
