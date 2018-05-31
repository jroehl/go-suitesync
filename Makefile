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

build: clean
	go build -o suitesync suitesync.go

clean:
	rm -rf suitesync

# docker-build:
# 	docker build -t jroehl/go-suitesync .

# docker-run:
# 	docker run --rm -it --name go-suitesync -p 5000:5000 jroehl/go-suitesync

# docker-push:
# 	docker push jroehl/go-suitesync

# test-unit:
# 	go test ./...
