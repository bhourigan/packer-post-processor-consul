package main

import (
	"github.com/bhourigan/post-processor-consul/consul"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(new(consul.PostProcessor))
	server.Serve()
}
