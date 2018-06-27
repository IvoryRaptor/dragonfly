package mq

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"fmt"
	"os"
	"log"
	"github.com/IvoryRaptor/dragonfly"
)

type Kafka struct {
	kernel IArrive
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func (k * Kafka)KafkaPublish(topic string,partition int32, actor []byte,payload []byte) error {
	deliveryChan := make(chan kafka.Event)
	err := k.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
			Key:            actor,
			Value:          payload,
		},
		deliveryChan)
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	if err != nil {
		fmt.Println(topic, err.Error())
		return err
	}
	return nil
}

func (k * Kafka)KafkaConfig(kernel dragonfly.IKernel, config map[interface{}]interface{}) error{
	k.kernel = kernel.(IArrive)
	var err error = nil
	host := fmt.Sprintf("%s:%d",config["host"],config["port"])
	k.producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": host})
	if err != nil {
		return err
	}
	t := kafka.ConfigMap{
		"bootstrap.servers":    host,
		"group.id":             "PostOffice",
		"session.timeout.ms":   6000,
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"}}
	k.consumer, err = kafka.NewConsumer(&t)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	return nil
}

func (k * Kafka)Start() error {
	log.Printf("mq start %s_%s",k.kernel.Get("matrix"), k.kernel.Get("angler"))
	err := k.consumer.SubscribeTopics([]string{fmt.Sprintf("%s_%s", k.kernel.Get("matrix"), k.kernel.Get("angler"))}, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	go func() {
		for true {
			ev,err := k.consumer.ReadMessage(-1)
			if err != nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				k.kernel.Arrive(e.Value)
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}()
	//mq.consumer.Close()
	return nil
}

func (k * Kafka)Stop(){
	k.kernel.RemoveService(k)
}
