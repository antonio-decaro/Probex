package main

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nuclio/nuclio-sdk-go"
)

const USERNAME_ENV string = "MQTT_USERNAME" // "guest"
const PASSWORD_ENV string = "MQTT_PASSWORD" // "guest"
const IP_ENV string = "MQTT_BROKER_IP"      // "192.168.1.20"
const PORT_ENV string = "MQTT_PORT"         // "1883"
const ID string = "go_mqtt_logger"
const TOPIC string = "iot/mqtt/logger"
const QOS byte = 0

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	username := os.Getenv(USERNAME_ENV)
	password := os.Getenv(PASSWORD_ENV)
	ip := os.Getenv(IP_ENV)
	port := os.Getenv(PORT_ENV)

	body := string(event.GetBody())
	context.Logger.Info(body)

	client := getClient(
		ID,
		ip,
		port,
		username,
		password,
	)

	publish(client, QOS, TOPIC, body)

	return nuclio.Response{}, nil
}

func publish(client mqtt.Client, qos byte, topic string, data interface{}) {
	token := client.Publish(topic, qos, false, data)
	token.Wait()
}

func getClient(id, broker, port, username, password string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", broker, port))
	opts.SetClientID(id)
	opts.SetUsername(username)
	opts.SetPassword(password)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}
