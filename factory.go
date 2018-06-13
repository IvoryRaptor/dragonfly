package dragonfly

type RedisFactory struct {
}

func (f *RedisFactory)GetName() string{
	return "redis"
}

func (f *RedisFactory)Create(kernel IKernel,config map[interface {}]interface{}) (IService,error) {
	result := Redis{}
	result.Config(kernel, config)
	return &result, nil
}
