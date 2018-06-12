package test

import (
	"github.com/IvoryRaptor/dragonfly"
)

type ServiceFactory struct {
}

func (f * ServiceFactory)GetName() string{
	return "test"
}

func (f * ServiceFactory)Create(kernel dragonfly.IKernel,config map[interface {}]interface{}) dragonfly.IService {
	r := TestService{}
	r.Config(kernel, config)
	return &r
}
