package dragonfly


type IService interface {
	GetName() string
	Config(kernel IKernel, config interface{}) error
	Start() error
	Stop()
}
