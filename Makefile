
all: compile

## compile: Compile go program
compile: go-clean go-get go-build

## install: Install dependencies (go get)
install: go-get

go-build:
	go build

go-generate:
	go generate

go-get:
	go get

go-install:
	go install

go-clean:
	rm -f gokmp
	go clean
