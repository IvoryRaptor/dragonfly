package main

import (
	"github.com/IvoryRaptor/dragonfly/test"
	"github.com/IvoryRaptor/dragonfly"
	"log"
)

func main() {
	k := test.TestKernel{}
	k.New("test")
	err := dragonfly.Builder(
		&k,
		[]dragonfly.IServiceFactory{&test.ServiceFactory{}})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = k.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	k.WaitStop()
}
