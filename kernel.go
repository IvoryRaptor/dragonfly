package dragonfly

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type IKernel interface {
	GetConfig() (map[string]interface{}, error)
	Get(name string) interface{}
	Set(name string, value interface{})
	RemoveService(service interface{})
	AddService(name string, service IService)
	GetService(name string) IService
}

type Kernel struct {
	Name       string
	wait       sync.WaitGroup
	services   map[string]IService
	signalChan chan os.Signal
	parameter  map[string]interface{}
}

func (k *Kernel) New(name string) {
	k.Name = name
	k.services = map[string]IService{}
	k.parameter = map[string]interface{}{}
	k.signalChan = make(chan os.Signal, 1)
}

func (k *Kernel) GetConfig() (map[string]interface{}, error) {
	config := map[string]interface{}{}
	data, err := ioutil.ReadFile(fmt.Sprintf("./config/%s/config.yaml", k.Name))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, config)
	return config, err
}

func (k *Kernel) Get(name string) interface{} {
	return k.parameter[name]
}

func (k *Kernel) Set(name string, value interface{}) {
	k.parameter[name] = value
}

func (k *Kernel) AddService(name string, service IService) {
	k.services[name] = service
	k.wait.Add(1)
}

func (k *Kernel) RemoveService(service interface{}) {
	var name string = ""
	for n, s := range k.services {
		if service == s {
			name = n
			break
		}
	}
	delete(k.services, name)
	k.wait.Done()
}

func (k *Kernel) GetService(name string) IService {
	return k.services[name]
}

func (k *Kernel) Start() error {
	log.Printf("%s Start\n", k.Name)
	for name, service := range k.services {
		log.Printf("Service [%s] Start\n", name)
		err := service.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Kernel) Stop() {
	k.signalChan <- syscall.SIGINT
}

func (k *Kernel) WaitStop() {
	signal.Notify(k.signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-k.signalChan
	for name, service := range k.services {
		service.Stop()
		log.Printf("Service [%s] Stop\n", name)
	}
	log.Printf("Wait Stop\n")
	k.wait.Wait()
	os.Exit(0)
}
