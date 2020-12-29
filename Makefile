
VERSION = $(shell git describe --tags --always)

all: clean deps binaries

binaries:
	GOOS=windows GOARCH=amd64 go build -o bin/gokmp-windows-$(VERSION).exe
	GOOS=darwin GOARCH=amd64 go build -o bin/gokmp-darwin-$(VERSION)
	GOOS=linux GOARCH=amd64 go build -o bin/gokmp-linux-$(VERSION)

deps:
	go generate
	go get

clean:
	rm -rf gokmp-* bin/ cmd/flickr_config.go
	go clean
