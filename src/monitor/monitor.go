package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

const USERNAME_ENV string = "MQTT_USERNAME"
const PASSWORD_ENV string = "MQTT_PASSWORD"
const IP_ENV string = "MQTT_BROKER_IP"
const PORT_ENV string = "PORT_ENV"
const LOG_QUEUE_NAME string = "iot/monitor"

var (
	planetinfo []map[string]interface{}
	headers    = []string{"Name"}
)

func main() {
	fmt.Println("[*] Starting monitor...")
	startListening()
}

func startListening() {
	err := godotenv.Load("../../.env")
	failOnError(err, "Failed to read environment variables")

	username := os.Getenv(USERNAME_ENV)
	password := os.Getenv(PASSWORD_ENV)
	ip := os.Getenv(IP_ENV)
	port := os.Getenv(PORT_ENV)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, ip, port))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		LOG_QUEUE_NAME, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			updateData(d.Body)
		}
	}()

	forever := make(chan bool)
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func updateData(msg []byte) {

	var data map[string]interface{}

	err := json.Unmarshal(msg, &data)
	if err == nil {
		fmt.Printf("%+v\n", data)

	} else {
		fmt.Println(err)
	}
}
