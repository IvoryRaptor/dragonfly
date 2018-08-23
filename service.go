package dragonfly

type IService interface {
	Config(kernel IKernel, config map[interface{}]interface{}) error
	Start() error
	Stop()
}

type IServiceFactory interface {
	GetName() string
	Create(kernel IKernel, config map[interface{}]interface{}) (IService, error)
}

func Builder(kernel IKernel, factories []IServiceFactory) error {
	config, err := kernel.GetConfig()
	if err != nil {
		return err
	}
	for _, factory := range factories {
		var service IService
		println("[" + factory.GetName() + "]")
		service, err = factory.Create(kernel, config[factory.GetName()].(map[interface{}]interface{}))
		kernel.AddService(factory.GetName(), service)
	}
	return nil
}
