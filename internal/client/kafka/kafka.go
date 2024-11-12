package kafka

type Producer interface {
	Send(topicName string, value string) error
	Close() error
}
