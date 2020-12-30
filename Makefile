
VERSION = $(shell git describe --tags --always)
IMAGE = sorend/gokmp

all: clean deps binaries docker

binaries:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/gokmp-windows-$(VERSION).exe
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/gokmp-darwin-$(VERSION)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/gokmp-linux-$(VERSION)

deps:
	go generate
	go get

clean:
	rm -rf gokmp-* bin/ cmd/flickr_config.go
	go clean

docker:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(VERSION) .
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
