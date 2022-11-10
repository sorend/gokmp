
VERSION = $(shell git describe --tags --always)
IMAGE = sorend/gokmp

default: help

# Builds gokmp binaries as well as docker image (requires docker setup)
all: build docker

# Builds gokmp binaries
build: deps
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/gokmp-windows-$(VERSION).exe
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/gokmp-darwin-$(VERSION)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/gokmp-linux-$(VERSION)

# Generates go source and get dependencies
deps: clean
	go generate
	go get

# Remove binaries and generated source
clean:
	rm -rf gokmp-* bin/ cmd/flickr_config.go
	go clean

# Create $(IMAGE) image
docker: build
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(VERSION) .
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest

# Push docker image
docker-deploy:
	# echo "$(DOCKER_PASSWORD)" | docker login -u $(DOCKER_USERNAME) --password-stdin
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

# Show this help prompt.
help:
	@ echo
	@ echo '  Usage:'
	@ echo ''
	@ echo '    make <target> [flags...]'
	@ echo ''
	@ echo '  Targets:'
	@ echo ''
	@ awk '/^#/{ comment = substr($$0,3) } comment && /^[a-zA-Z][a-zA-Z0-9_-]+ ?:/{ print "   ", $$1, comment }' $(MAKEFILE_LIST) | column -t -s ':' | sort
	@ echo ''
