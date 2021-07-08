package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ProbeData struct {
	Name         string
	Coordinate   [2]float64
	Distance     float64
	StarDistance float64 // astronomic units
	StarType     string
	Mass         float64
	Radius       float64
	Humidity     float32
	Temperature  float32
	Wind         float32
	ProbeId      uint64
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

	var data ProbeData
	err := json.Unmarshal(msg.Payload(), &data)
	if err != nil {
		fmt.Printf("[ERROR] Payload: %s\n", string(msg.Payload()))
		return
	}

	fmt.Printf("[*] Sending probe n %d on planet: %+v\n", id, data)
	time.Sleep(3 * time.Second)

	collectPlanetData(&data, id)
	send, _ := json.Marshal(data)

	fmt.Printf("[*] Probe %d retrived those information about the planet %s: %s\n", id, data.Name, string(send))

	if token := client.Publish("iot/probe/receiver", 1, false, send); token.Wait() && token.Error() != nil {
		fmt.Printf("[ERROR] %s\n", token.Error())
	}
}

func collectPlanetData(data *ProbeData, id uint64) {

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	data.Humidity = r.Float32() + float32(r.Intn(10))
	data.Temperature = r.Float32() + float32(r.Intn(10))
	data.Wind = r.Float32() + float32(r.Intn(10))
}
