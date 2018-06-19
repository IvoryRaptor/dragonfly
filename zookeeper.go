package dragonfly

import (
	"github.com/samuel/go-zookeeper/zk"
	"fmt"
	"sync"
	"time"
)

type ZookeeperFactory struct {
}

func (f *ZookeeperFactory) GetName() string {
	return "zookeeper"
}

func (f *ZookeeperFactory) Create(kernel IKernel, config map[interface{}]interface{}) (IService, error) {
	result := Zookeeper{}
	result.Config(kernel, config)
	return &result, nil
}

type ZkNode struct {
	path     string
	stopChan chan bool
	level    int
	childes  sync.Map
}

func (n *ZkNode) GetChildes() map[string]*ZkNode {
	result := map[string]*ZkNode{}
	n.childes.Range(func(k, v interface{}) bool {
		result[k.(string)] = v.(*ZkNode)
		return true
	})
	return result
}

func (n *ZkNode) GetKeys() []string {
	result := []string{}
	n.childes.Range(func(k, v interface{}) bool {
		result = append(result,k.(string))
		return true
	})
	return result
}

func (n *ZkNode) Watch(conn *zk.Conn, path string, level int) {
	n.childes = sync.Map{}
	go func() {
		for {
			keys, _, childCh, _ := conn.ChildrenW(path)
			newMap := sync.Map{}
			for _, node := range keys {
				iotnn, ok := n.childes.Load(node)
				if ok {
					newMap.Store(node, iotnn)
					n.childes.Delete(node)
				} else {
					iotnn := &ZkNode{}
					if level > 0 {
						iotnn.Watch(conn, path+"/"+node, level-1)
					}
					newMap.Store(node, iotnn)
				}
			}
			n.childes.Range(func(k, v interface{}) bool {
				v.(*ZkNode).StopWatch()
				return true
			})
			n.childes = newMap
			select {
			case ev := <-childCh:
				if ev.Err != nil {
					fmt.Printf("Child watcher error: %+v\n", ev.Err)
					return
				}
			case <-n.stopChan:
				return
			}
		}
	}()
}

func (n *ZkNode) StopWatch() {
	n.stopChan <- true
}

type Zookeeper struct {
	ZkNode
	Kernel IKernel
	url    string
	conn   *zk.Conn
	level  int
	path   string
}

func (z *Zookeeper) GetConn() *zk.Conn {
	return z.conn
}

func (z *Zookeeper) Config(kernel IKernel, config map[interface{}]interface{}) error {
	z.Kernel = kernel
	z.url = fmt.Sprintf("%s:%d", config["host"].(string), config["port"].(int))
	z.path = config["path"].(string)
	z.level = config["level"].(int)
	return nil
}

func (z *Zookeeper) Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	return z.conn.Create(path, data, flags, acl)
}

func (z *Zookeeper) Start() error {
	var err error
	z.conn, _, err = zk.Connect([]string{z.url}, time.Second*10)
	if err != nil {
		return err
	}
	z.Watch(z.conn, z.path, z.level)
	return nil
}

func (z *Zookeeper) Stop() {
	z.StopWatch()
	z.Kernel.RemoveService(z)
}
