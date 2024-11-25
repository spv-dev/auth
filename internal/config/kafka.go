package config

import (
	"errors"
	"os"
	"strings"

	"github.com/IBM/sarama"

	serviceerror "github.com/spv-dev/auth/internal/service_error"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"

	saramaRetryMax        = 5
	saramaReturnSuccesses = true
)

// KafkaProducerConfig интерфейс для работы с Kafka
type KafkaProducerConfig interface {
	Brokers() []string
	GroupID() string
	Config() *sarama.Config
}

type kafkaProducerConfig struct {
	brokers []string
	groupID string
}

// NewKafkaProducerConfig получение конфигурации для подключения к Kafka
func NewKafkaProducerConfig() (*kafkaProducerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New(serviceerror.KafkaBrokersNotFound)
	}

	brokers := strings.Split(brokersStr, ",")

	groupID := os.Getenv(groupIDEnvName)
	if len(groupID) == 0 {
		return nil, errors.New(serviceerror.KafkaGroupIDNotFound)
	}

	return &kafkaProducerConfig{
		brokers: brokers,
		groupID: groupID,
	}, nil
}

func (cfg *kafkaProducerConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaProducerConfig) GroupID() string {
	return cfg.groupID
}

func (cfg *kafkaProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = saramaRetryMax
	config.Producer.Return.Successes = saramaReturnSuccesses

	return config
}
