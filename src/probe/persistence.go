package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type Persistence struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func InitPersistence() (*Persistence, error) {
	username := os.Getenv(USERNAME_ENV)
	password := os.Getenv(PASSWORD_ENV)
	ip := os.Getenv(IP_ENV)
	port := os.Getenv(PORT_ENV)

	ret := new(Persistence)

	var err error
	ret.conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, ip, port))
	if err != nil {
		return ret, err
	}

	ret.ch, err = ret.conn.Channel()
	if err != nil {
		return ret, err
	}

	ret.queue, err = ret.ch.QueueDeclare(
		MONITOR_QUEUE_NAME, // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func (p *Persistence) PersistProbeData(data ProbeData) error {

	msg, _ := json.Marshal(data)

	err := p.ch.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
	return err
}
