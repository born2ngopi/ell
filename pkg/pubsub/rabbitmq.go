package pubsub

import (
	"fmt"
	"log"

	"github.com/born2ngopi/eel/pkg/memcache"
	"github.com/streadway/amqp"
)

type rabbitmq struct {
	conn *amqp.Connection
}

type RabbitOption struct {
	Username string
	Password string
	Host     string
	Port     string
}

func NewRabbit(opt RabbitOption) (Pubsub, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", opt.Username, opt.Password, opt.Host, opt.Port))
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	return &rabbitmq{
		conn: conn,
	}, nil
}

func (r *rabbitmq) Subscribe(topics []string) {

	ch, err := r.conn.Channel()
	if err != nil {
		log.Fatal("failed to open a channel : ", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal("failed to declare a queue : ", err)
	}

	for _, topic := range topics {
		err = ch.QueueBind(
			q.Name,      // queue name
			topic,       // routing key
			"amq.topic", // exchange
			false,
			nil)
		if err != nil {
			log.Fatal("failed to bind a queue : ", err)
		}
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal("failed to consume a queue : ", err)
	}

	for msg := range msgs {
		memcache.Update(msg.RoutingKey, string(msg.Body))
	}

}

func (r *rabbitmq) Publish(topic string, message []byte) error {

	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = ch.Publish(
		"amq.topic", // exchange
		topic,       // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		return err
	}
	return nil
}
