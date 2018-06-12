package test

import "github.com/IvoryRaptor/dragonfly"

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
