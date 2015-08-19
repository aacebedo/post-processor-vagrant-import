package main

import (
        "github.com/mitchellh/packer/packer/plugin"
        "vagrant-import"
       )

func main() {
	server, err := plugin.Server()
    if err != nil {
    	panic(err)
    }
	server.RegisterPostProcessor(new(vagrantimport.PostProcessor))
	server.Serve()
}
