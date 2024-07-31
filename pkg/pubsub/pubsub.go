package pubsub

type Pubsub interface {
	Subscribe(topic []string)
	Publish(topic string, message []byte) error
}
