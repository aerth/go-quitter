# aerth [https://isupon.us]
# https://github.com/aerth

NAME=go-quitter
XTRA?=
VERSION=0.0.9
RELEASE:=v${VERSION}.$(shell git rev-parse --verify --short HEAD)
PREFIX=/usr/local/bin
LDFLAGS=-s
GOCMD=./cmd/go-quitter
CGO_ENABLED?=0

all:
	@echo ""
	@echo "building go-quitter command line GNU Social client to"
	@echo "		\"bin/${NAME}-${RELEASE}\""
	@echo
	@mkdir -p bin
	@go build -v -x ${XTRA} -o bin/${NAME} -ldflags '${LDFLAGS}' ${GOCMD}

deps:
	go get -d ./...

test:
	go test -v

install: bin/go-quitter
#	sudo mv bin/${NAME}-${RELEASE} ${PREFIX}/${NAME}
	sudo mv bin/${NAME} ${PREFIX}/${NAME}

update: upgrade deps all install

upgrade:
	git pull origin master

cross:
	mkdir -p bin
	GOOS=windows GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-WIN32.exe -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=windows GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-WIN64.exe -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=darwin GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-OSX-x86 -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=darwin GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-OSX-amd64 -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=linux GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-linux-amd64 -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=linux GOARCH=arm go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-linux-arm -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=linux GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-linux-x86 -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=freebsd GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-freebsd-amd64 -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=freebsd GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-freebsd-x86 -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=openbsd GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-openbsd-amd64 -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=openbsd GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-openbsd-x86 -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=netbsd GOARCH=amd64 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-netbsd-amd64 -ldflags='${LDFLAGS}' ${GOCMD}
	GOOS=netbsd GOARCH=386 go build -v -x ${XTRA} -o bin/${NAME}-${RELEASE}-netbsd-x86 -ldflags='${LDFLAGS}' ${GOCMD}

cui:
	XTRA='-tags cui' make
