package main

import (
	"github.com/IvoryRaptor/dragonfly"
	"github.com/IvoryRaptor/dragonfly/test"
	"log"
)

func main() {
	k := test.TestKernel{}
	k.New("test", k.SetFields)
	err := dragonfly.Builder(
		&k,
		[]dragonfly.IServiceFactory{
			&test.Factory{},
			&dragonfly.ZookeeperFactory{},
		})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = k.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	k.WaitStop()
}
