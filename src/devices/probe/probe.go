package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ProbeData struct {
	Name        string
	Humidity    float32
	Temperature float32
	Wind        float32
}

type ProbeInfo struct {
	Name       string
	Coordinate [2]float64
	Distance   float64
}

var probeid uint64

func main() {
	fmt.Println("[*] Starting Probe simulation device.")
	defer fmt.Println("[.] Terminating Probe simulation device")

	opts := mqtt.NewClientOptions()

	const (
		username = "guest"
		password = "guest"
		ip       = "localhost"
	)

	opts.AddBroker(fmt.Sprintf("tcp://%s:1883", ip))
	opts.SetClientID("telescope")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(probe)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(0)

	if token := client.Subscribe("iot/probe", 2, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	forever := make(chan bool)
	<-forever
}

func probe(client mqtt.Client, msg mqtt.Message) {
	id := atomic.AddUint64(&probeid, 1)

	var data ProbeInfo
	err := json.Unmarshal(msg.Payload(), &data)
	if err != nil {
		fmt.Printf("[ERROR] Payload: %s\n", string(msg.Payload()))
		return
	}

	fmt.Printf("[*] Sending probe n %d on planet: %+v\n", id, data)

	// TODO collezionare dati della sonda
	planetData := collectPlanetData(data.Name)
	jsonVal, _ := json.Marshal(planetData)

	time.Sleep(3 * time.Second)

	fmt.Printf("[*] Probe %d retrived those information about the planet %s: %s\n", id, data.Name, string(jsonVal))

	if token := client.Publish("iot/probe/receiver", 1, false, jsonVal); token.Wait() && token.Error() != nil {
		fmt.Printf("[ERROR] %s\n", token.Error())
	}
}

func collectPlanetData(name string) *ProbeData {
	return &ProbeData{
		Name:        name,
		Humidity:    4, // random
		Temperature: 3, // random
		Wind:        1, // random
	}
}
