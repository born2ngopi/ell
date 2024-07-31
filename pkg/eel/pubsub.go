package eel

import (
	"errors"
	"log"

	"github.com/born2ngopi/eel/pkg/client"
	"github.com/born2ngopi/eel/pkg/memcache"
	"github.com/born2ngopi/eel/pkg/pubsub"
	"github.com/robfig/cron/v3"
)

type WatchOption struct {
	// standard interval can be read from https://en.wikipedia.org/wiki/Cron
	Interval string
	// Driver is the driver to use for the pubsub system.
	// support: nsq, rabbitmq
	Driver string
	// Host is the host of the pubsub system.
	Host string
	// Port is the port of the pubsub system.
	Port string
	// Username is the username of the pubsub system.
	Username string
	// Password is the password of the pubsub system.
	Password string
}

const (
	RABBITMQ_DRIVER = "rabbitmq"
)

// StartPubSub is function for watch update if any new token from key
// parameter opt is the PubSubOption for the pubsub system.
// and keys optional, if keys is empty, then watch all keys.
// and if keys is not empty, then watch only keys.
func Watch(opt WatchOption, key []string) error {

	var (
		broker pubsub.Pubsub
		err    error
	)

	// todo
	if opt.Driver == RABBITMQ_DRIVER {

		broker, err = pubsub.NewRabbit(pubsub.RabbitOption{
			Username: opt.Username,
			Password: opt.Password,
			Host:     opt.Host,
			Port:     opt.Port,
		})
		if err != nil {
			return err
		}
		// todo
	} else {
		return errors.New("driver not supported")
	}

	log.Println("Start PubSub to watch update token")
	go broker.Subscribe(key)

	if opt.Interval != "" {
		c := cron.New()

		_, err = c.AddFunc(opt.Interval, func() {
			for _, k := range key {
				token, err := client.Do(e.token, e.host, k)
				if err != nil {
					log.Printf("failed to get token from key: %s, with error %s ", k, err)
				} else {
					memcache.Set(e.token, token)
				}
			}
		})
		if err != nil {
			return err
		}

		c.Start()
	}

	return nil
}
