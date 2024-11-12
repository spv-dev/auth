package producer

import (
	"log"

	"github.com/IBM/sarama"
)

type producer struct {
	sender sarama.SyncProducer
}

func NewProducer(sender sarama.SyncProducer) *producer {
	return &producer{
		sender: sender,
	}
}

func (p *producer) Send(topicName string, value string) error {
	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(value),
	}

	log.Printf("try to send message to kafka: %#v", msg)

	partition, offset, err := p.sender.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message in Kafka: %v\n", err.Error())
		return err
	}

	log.Printf("message sent to partition %d with offset %d\n", partition, offset)
	return nil
}

func (p *producer) Close() error {
	return p.sender.Close()
}
