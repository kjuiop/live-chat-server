package kafka

type Client interface {
	Subscribe(topic string) error
}
