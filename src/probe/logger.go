package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

const USERNAME_ENV string = "MQTT_USERNAME"
const PASSWORD_ENV string = "MQTT_PASSWORD"
const IP_ENV string = "MQTT_BROKER_IP"
const PORT_ENV string = "PORT_ENV"
const LOG_QUEUE_NAME string = "iot/logs"

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

func (logger *Logger) log(msg string) error {
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

func (logger *Logger) Info(msg string) error {
	msg = "[INFO] " + msg
	return logger.log(msg)
}

func (logger *Logger) Debug(msg string) error {
	msg = "[DEBUG] " + msg
	return logger.log(msg)
}

func (logger *Logger) Error(msg string) error {
	msg = "[ERROR] " + msg
	return logger.log(msg)
}

func (logger *Logger) Warning(msg string) error {
	msg = "[WARNING] " + msg
	return logger.log(msg)
}

func (logger *Logger) Close() {
	defer logger.ch.Close()
	defer logger.conn.Close()
}
