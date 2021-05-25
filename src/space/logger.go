package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type Logger struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func InitLogger() (*Logger, error) {
	username := os.Getenv(USERNAME_ENV)
	password := os.Getenv(PASSWORD_ENV)
	ip := os.Getenv(IP_ENV)
	port := os.Getenv(PORT_ENV)

	logger := new(Logger)

	var err error
	logger.conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, ip, port))
	if err != nil {
		return logger, err
	}

	logger.ch, err = logger.conn.Channel()
	if err != nil {
		return logger, err
	}

	logger.queue, err = logger.ch.QueueDeclare(
		LOG_QUEUE_NAME, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return logger, err
	}

	return logger, nil
}

func (logger *Logger) Log(msg string) error {
	err := logger.ch.Publish(
		"",                // exchange
		logger.queue.Name, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	return err
}

func (logger *Logger) Close() {
	defer logger.ch.Close()
	defer logger.conn.Close()
}
