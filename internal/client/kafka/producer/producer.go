package producer

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"

	serviceerror "github.com/spv-dev/auth/internal/service_error"
)

type producer struct {
	sender sarama.SyncProducer
}

// NewProducer получение объекта для взаимодействия с Kafka
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

	partition, offset, err := p.sender.SendMessage(msg)
	if err != nil {
		return fmt.Errorf(serviceerror.FailedToSendMessageInKafka, err.Error())
	}

	log.Printf("message sent to partition %d with offset %d\n", partition, offset)
	return nil
}

func (p *producer) Close() error {
	return p.sender.Close()
}
