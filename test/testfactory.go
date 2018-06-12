package test

import (
	"github.com/IvoryRaptor/dragonfly"
)

type Factory struct {
}

func (f *Factory)GetName() string{
	return "test"
}

func (f *Factory)Create(kernel dragonfly.IKernel,config map[interface {}]interface{}) (dragonfly.IService,error) {
	r := TestService{}
	err := r.Config(kernel, config)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
