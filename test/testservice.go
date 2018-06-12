package test

import (
	"github.com/IvoryRaptor/dragonfly"
	"time"
	"log"
)

type TestService struct {
	kernel dragonfly.IKernel
	name string
	run bool
}

func (t * TestService)GetName()string{
	return "test"
}

func (t * TestService)Config(kernel dragonfly.IKernel, config map[interface {}]interface{}) error {
	t.kernel = kernel
	t.name = config["name"].(string)
	return nil
}

func (t * TestService)Start() error {
	t.run = true
	go func() {
		for t.run  {
			println(t.name)
			time.Sleep(time.Second)
		}
		t.kernel.RemoveService(t)
	}()
	return nil
}

func (t * TestService)Stop() {
	log.Printf("Stop [%s] Service", t.name)
	t.run = false
}

