package main

//go:generate generate/flickr_config.sh cmd/flickr_config.go

import (
	"github.com/sorend/gokmp/cmd"
)


func main() {
	cmd.Execute()
}

