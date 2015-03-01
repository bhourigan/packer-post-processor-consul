package main

import (
	"bhourigan/packer-post-processor-ami/ami"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(new(ami.AMIPostProcessor))
	server.Serve()
}
