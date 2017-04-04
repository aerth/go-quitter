# aerth [https://isupon.us]
# https://github.com/aerth

NAME=go-quitter
XTRA?=
VERSION=0.0.9
#RELEASE?=v${VERSION}.$(shell git rev-parse --verify --short HEAD)
RELEASE:=v${VERSION}.$(shell git rev-parse --verify --short HEAD)
PREFIX=/usr/local/bin
DATE=$(shell date -u)
HOSTNAME=$(shell hostname)
LDFLAGS=-s -w -X main.release=$(RELEASE) -X 'main.buildinfo=$(DATE) ($(HOSTNAME))'
GOCMD=./cmd/go-quitter
CGO_ENABLED?=0

all: build

help:

	@echo "make" - To build the current go-quitter tree
	@echo "make cui" - Build go-quitter with Console User Interface (experimental)
	@echo "make update - To update the source code, build, and install to /usr/local/bin/"
	@echo "make install" mv bin/go-quitter $PREFIX/
	@echo "make test" To run go-quitter library test functions

build:
	@echo ""
	@echo "building go-quitter command line GNU Social client to"
	@echo "		\"bin/${NAME}\""
	@echo
	@mkdir -p bin
	go build -v -x ${XTRA} -o bin/${NAME} -ldflags="${LDFLAGS}" ${GOCMD}


test:
	go test -v ./...

install: bin/go-quitter
	mv bin/${NAME} ${PREFIX}/${NAME}
	chmod 755 ${PREFIX}/${NAME}
	@echo installed as ${PREFIX}/${NAME}
	@rmdir bin 2> /dev/null | true

update: upgrade deps build
	su -c 'make install'

upgrade:
	git pull origin master

deps:
		go get -v -d ./...

cross:
	mkdir -p bin
	GOOS=windows GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-WIN32.exe -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=windows GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-WIN64.exe -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=darwin GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-OSX-x86 -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=darwin GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-OSX-amd64 -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=linux GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-linux-amd64 -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=linux GOARCH=arm go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-linux-arm -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=linux GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-linux-x86 -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=freebsd GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-freebsd-amd64 -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=freebsd GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-freebsd-x86 -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=openbsd GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-openbsd-amd64 -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=openbsd GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-openbsd-x86 -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=netbsd GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-netbsd-amd64 -ldflags="${LDFLAGS}" ${GOCMD}
	GOOS=netbsd GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-netbsd-x86 -ldflags="${LDFLAGS}" ${GOCMD}

cui:
	XTRA='-tags cui' make
