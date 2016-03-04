NAME=go-quitter
RELEASE:=$(shell git rev-parse --verify --short HEAD)
USER=aerth
# First attempt at Makefile
all:	
	go fmt
	go vet
	cd cmd/go-quitter && go build -v -o ../../go-quitter

install:
	mv go-quitter /usr/local/bin/go-quitter
update:
	git pull origin master
	mv go-quitter /usr/local/bin/go-quitter
