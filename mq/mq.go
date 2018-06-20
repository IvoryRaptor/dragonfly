package mq

import (
	"github.com/IvoryRaptor/dragonfly"
)

type IArrive interface {
	dragonfly.IKernel
	Arrive([]byte)
}

type IMQ interface {
	dragonfly.IService
	Publish(topic string, actor []byte, payload []byte) error
}
