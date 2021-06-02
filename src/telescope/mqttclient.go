package main

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTClient(id string) mqtt.Client {
	opts := mqtt.NewClientOptions()

	username := os.Getenv(USERNAME_ENV)
	password := os.Getenv(PASSWORD_ENV)
	ip := os.Getenv(IP_ENV)

	opts.AddBroker(fmt.Sprintf("tcp://%s:1883", ip))
	opts.SetClientID(id)
	opts.SetUsername(username)
	opts.SetPassword(password)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}
