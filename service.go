package dragonfly


type IService interface {
	GetName() string
	Config(kernel IKernel, config map[interface{}]interface{}) error
	Start() error
	Stop()
}

type IServiceFactory interface {
	GetName() string
	Create(kernel IKernel,config map[interface {}]interface{}) IService
}

func Builder(kernel IKernel, factories []IServiceFactory)error {
	config,err := kernel.GetConfig()
	if err != nil {
		return err
	}
	for _, factory := range factories {
		println("123")
		service := factory.Create(kernel, config[factory.GetName()].(map[interface {}]interface {}))
		kernel.AddService(service)
	}
	return nil
}