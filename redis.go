package dragonfly

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"sync"
)

type Redis struct {
	Mutex  sync.Mutex
	Conn   redis.Conn
	url    string
	kernel IKernel
}

func (r *Redis) Config(kernel IKernel, config map[interface{}]interface{}) error {
	r.kernel = kernel
	r.url = fmt.Sprintf("%s:%d", config["host"], config["port"])
	return nil
}

func (r *Redis) Start() error {
	r.Mutex = sync.Mutex{}
	log.Printf("Redis %s", r.url)
	c, err := redis.Dial("tcp", r.url)
	if err != nil {
		return err
	}
	r.Conn = c
	return nil
}

func (r *Redis) Do(commandName string, args ...interface{}) (interface{}, error) {
	r.Mutex.Lock()
	reply, err := r.Conn.Do(commandName, args...)
	r.Mutex.Unlock()
	return reply, err
}

func (r *Redis) Stop() {
	r.Conn.Close()
	r.kernel.RemoveService(r)
}

type RedisFactory struct {
}

func (f *RedisFactory) GetName() string {
	return "redis"
}

func (f *RedisFactory) Create(kernel IKernel, config map[interface{}]interface{}) (IService, error) {
	result := Redis{}
	result.Config(kernel, config)
	return &result, nil
}

var Singleton = RedisFactory{}
