package main

import (
	"github.com/IvoryRaptor/dragonfly"
	"log"
	"time"
)

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

type TestConfig struct {
	Name string `yaml:"name"`
}

type Config struct {
	Test TestConfig `yaml:"test"`
}

type TestKernel struct{
	dragonfly.Kernel
}

func (t * TestKernel)Config()error {
	var config Config
	err := t.InitializeConfig(&config)
	if err != nil {
		return err
	}
	server := TestService{}
	server.Config(t, config.Test)
	t.AddService(&server)
	return err
}

func main() {
	t :=TestKernel{}
	t.Name = "test"
	t.Config()
	t.Start()
	t.WaitStop()
}
