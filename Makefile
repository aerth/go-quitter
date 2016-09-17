# aerth [https://isupon.us]
# https://github.com/aerth

NAME=go-quitter
VERSION=0.0.9
RELEASE:=v${VERSION}.$(shell git rev-parse --verify --short HEAD)
PREFIX=/usr/local/bin
all:
	@echo
	@echo "building go-quitter command line GNU Social client to"
	@echo "		\"bin/${NAME}-${RELEASE}\""
	@echo
	@mkdir -p bin
	@go build -v -o bin/${NAME}-${RELEASE} ./cmd/go-quitter

deps:
	go get -d ./...


install:
	sudo mv bin/${NAME}-${RELEASE} ${PREFIX}/${NAME}


update: upgrade deps all install

upgrade:
	git pull origin master

cross:
	mkdir -p bin
	GOOS=windows GOARCH=386 go build -v -o bin/${NAME}-${RELEASE}-WIN32.exe
	GOOS=windows GOARCH=amd64 go build -v -o bin/${NAME}-${RELEASE}-WIN64.exe
	GOOS=darwin GOARCH=386 go build -v -o bin/${NAME}-${RELEASE}-OSX-x86
	GOOS=darwin GOARCH=amd64 go build -v -o bin/${NAME}-${RELEASE}-OSX-amd64
	GOOS=linux GOARCH=amd64 go build -v -o bin/${NAME}-${RELEASE}-linux-amd64
	GOOS=linux GOARCH=386 go build -v -o bin/${NAME}-${RELEASE}-linux-x86
	GOOS=freebsd GOARCH=amd64 go build -v -o bin/${NAME}-${RELEASE}-freebsd-amd64
	GOOS=freebsd GOARCH=386 go build -v -o bin/${NAME}-${RELEASE}-freebsd-x86
	GOOS=openbsd GOARCH=amd64 go build -v -o bin/${NAME}-${RELEASE}-openbsd-amd64
	GOOS=openbsd GOARCH=386 go build -v -o bin/${NAME}-${RELEASE}-openbsd-x86
	GOOS=netbsd GOARCH=amd64 go build -v -o bin/${NAME}-${RELEASE}-netbsd-amd64
	GOOS=netbsd GOARCH=386 go build -v -o bin/${NAME}-${RELEASE}-netbsd-x86
