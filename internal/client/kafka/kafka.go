package kafka

// Producer интерфейс для отправки сообщений в Kafka
type Producer interface {
	Send(topicName string, value string) error
	Close() error
}
