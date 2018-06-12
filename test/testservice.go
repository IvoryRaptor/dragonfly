package test

import (
	"github.com/IvoryRaptor/dragonfly"
	"time"
	"log"
)

type TestConfig struct {
	Name string `yaml:"name"`
}

type TestService struct {
	kernel dragonfly.IKernel
	name string
	run bool
}

func (t * TestService)GetName()string{
	return "test"
}

func (t * TestService)Config(kernel dragonfly.IKernel, config interface{}) error {
	t.kernel = kernel
	t.name = config.(TestConfig).Name
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
	log.Printf("Stop Service %s", t.name)
	t.run = false
}

