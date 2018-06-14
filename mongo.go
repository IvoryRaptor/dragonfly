package dragonfly

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"time"
	"log"
)

type Mongo struct {
	url string
	kernel IKernel
	session *mgo.Session
}

func (m *Mongo)GetSession() *mgo.Session {
	return m.session.Clone()
}

func (m *Mongo) Config(kernel IKernel, config map[interface {}]interface{}) error {
	m.kernel = kernel
	m.url = fmt.Sprintf("mongodb://%s:%d", config["host"].(string), config["port"].(int))
	return nil
}

func (m *Mongo) Start() error{
	log.Printf("MongoDB Url %s", m.url)
	session, err := mgo.DialWithTimeout(m.url, 10*time.Second)
	if err != nil {
		return err
	}
	m.session = session
	m.session.SetMode(mgo.Monotonic, true)
	m.session.SetPoolLimit(300)
	return nil
}


func (m *Mongo) Stop(){
	m.kernel.RemoveService(m)
}
