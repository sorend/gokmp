
VERSION = $(shell git describe --tags --always)

all: compile

## compile: Compile go program
compile: clean go-get generate binaries

## install: Install dependencies (go get)
install: go-get

binaries:
	GOOS=windows GOARCH=amd64 go build -o bin/gokmp-windows-$(VERSION).exe
	GOOS=darwin GOARCH=amd64 go build -o bin/gokmp-darwin-$(VERSION)
	GOOS=linux GOARCH=amd64 go build -o bin/gokmp-linux-$(VERSION)

generate:
	go generate

go-get:
	go get

clean:
	rm -rf gokmp-* bin/
	go clean
