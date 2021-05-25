package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nuclio/nuclio-sdk-go"
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

type ProbeData struct {
	Name        string
	Distance    float32
	Mass        int32
	Radius      float32
	Temperature float32
	Water       bool
}

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	logger, err := InitLogger()
	if err != nil {
		context.Logger.Error("Error: %s", err)
		panic(err)
	}
	defer logger.Close()

	// if we got the event from rabbit
	if event.GetTriggerInfo().GetClass() == "async" && event.GetTriggerInfo().GetKind() == "mqtt" {
		body := event.GetBody()
		logger.log("Body content: " + string(body))

		var p ProbeData
		json.Unmarshal(body, &p)

		context.Logger.Info(p)
		logger.log(fmt.Sprint(p))

		return nil, nil
	}

	return nuclio.Response{
		StatusCode:  http.StatusOK,
		ContentType: "application/json",
	}, nil
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

func (logger *Logger) Close() {
	defer logger.ch.Close()
	defer logger.conn.Close()
}