package test

import (
	"github.com/IvoryRaptor/dragonfly"
	"time"
	"log"
)

type TestKernel struct {
	dragonfly.Kernel
}

func (k * TestKernel)T() {
	println(123)
}

type TestService struct {
	kernel *TestKernel
	name string
	run bool
}

func (t * TestService)GetName()string{
	return "test"
}

func (t * TestService)Config(kernel dragonfly.IKernel, config map[interface {}]interface{}) error {
	t.kernel = kernel.(*TestKernel)
	t.name = config["name"].(string)
	return nil
}

func (t * TestService)Start() error {
	t.run = true
	go func() {
		for t.run  {
			println(t.name)
			t.kernel.T()
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

