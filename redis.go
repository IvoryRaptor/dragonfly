package dragonfly

import (
	"sync"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"log"
)

type Redis struct {
	Mutex sync.Mutex
	Conn  redis.Conn
	url   string
}

func (r *Redis) Config(kernel IKernel, config map[interface {}]interface{}) error {
	r.url = fmt.Sprintf("%s:%d", config["host"], config["port"])
	return nil
}

func (r *Redis) Start() error {
	r.Mutex = sync.Mutex{}
	log.Printf("Redis %s",r.url)
	c, err := redis.Dial("tcp", r.url)
	if err != nil {
		return err
	}
	r.Conn = c
	return nil
}

func (r *Redis) Do(commandName string, args ...interface{}) (interface{}, error) {
	r.Mutex.Lock()
	reply, err := r.Conn.Do(commandName, args)
	r.Mutex.Unlock()
	return reply, err
}

func (r *Redis) Stop(){
	r.Conn.Close()
}