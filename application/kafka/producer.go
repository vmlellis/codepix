package kafka

import (
	"fmt"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	Producer     *ckafka.Producer
	DeliveryChan chan ckafka.Event
}

func NewKafkaProducer(deliveryChan chan ckafka.Event) (*KafkaProducer, error) {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
	}
	p, err := ckafka.NewProducer(configMap)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer:     p,
		DeliveryChan: deliveryChan,
	}, nil
}

func (p *KafkaProducer) Publish(msg, topic string) error {
	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Value: []byte(msg),
	}
	err := p.Producer.Produce(message, p.DeliveryChan)
	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProducer) DeliveryReport() {
	for e := range p.DeliveryChan {
		switch ev := e.(type) {
		case *ckafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Delivery failed:", ev.TopicPartition)
			} else {
				fmt.Println("Delivered message to:", ev.TopicPartition)
			}
		}
	}
}
