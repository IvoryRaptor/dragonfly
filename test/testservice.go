package test

import (
	"github.com/IvoryRaptor/dragonfly"
	"log"
	"time"
)

type TestKernel struct {
	dragonfly.Kernel
	zookeeper *dragonfly.Zookeeper
}

func (t *TestKernel) T() {

}
func (t *TestKernel) New() error {
	t.NewKernel("test")
	err := dragonfly.Builder(
		t,
		[]dragonfly.IServiceFactory{
			&Factory{},
			&dragonfly.ZookeeperFactory{},
		})
	t.zookeeper = t.GetService("zookeeper").(*dragonfly.Zookeeper)
	return err
}

type TestService struct {
	kernel *TestKernel
	name   string
	run    bool
}

func (t *TestService) GetName() string {
	return "test"
}

func (t *TestService) Config(kernel dragonfly.IKernel, config map[interface{}]interface{}) error {
	t.kernel = kernel.(*TestKernel)
	t.name = config["name"].(string)
	return nil
}

func (t *TestService) Start() error {
	t.run = true
	go func() {
		for t.run {
			t.kernel.T()
			for k, v := range t.kernel.zookeeper.GetChildes() {
				println(k)
				for _, t := range v.GetKeys() {
					print("\t" + t)
				}
				println()
			}
			time.Sleep(3 * time.Second)
		}
	}()
	return nil
}

func (t *TestService) Stop() {
	log.Printf("Stop [%s] Service", t.name)
	t.kernel.RemoveService(t)
	t.run = false
}
