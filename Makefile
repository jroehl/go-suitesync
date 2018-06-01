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

# test-unit:
# 	go test ./...
