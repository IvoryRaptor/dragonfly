package test

import (
	"github.com/IvoryRaptor/dragonfly"
)

type ServiceFactory struct {
}

func (f * ServiceFactory)GetName() string{
	return "test"
}

func (f * ServiceFactory)Create(kernel dragonfly.IKernel,config map[interface {}]interface{}) (dragonfly.IService,error) {
	r := TestService{}
	err := r.Config(kernel, config)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
