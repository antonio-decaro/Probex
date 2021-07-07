package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type TelescopeData struct {
	Name         string
	Coordinate   [2]float64
	Distance     float64
	StarDistance float64 // astronomic units
	StarType     string
	Mass         float64
	Radius       float64
}

func main() {
	fmt.Println("[*] Starting Telescope simulation device")
	defer fmt.Println("[.] Terminating Telescope simulation device")

	var data TelescopeData
	var x, y, distance, stardistance, mass, radius string
	var xfloat, yfloat float64

	fmt.Println("[+] Insert Name: ")
	fmt.Scanln(&data.Name)
	fmt.Println("[+] Insert Coordinate (x y): ")
	fmt.Scanln(&x, &y)
	fmt.Println("[+] Insert Distance: ")
	fmt.Scanln(&distance)
	fmt.Println("[+] Insert StarDistance: ")
	fmt.Scanln(&stardistance)
	fmt.Println("[+] Insert StarType: ")
	fmt.Scanln(&data.StarType)
	fmt.Println("[+] Insert Mass: ")
	fmt.Scanln(&mass)
	fmt.Println("[+] Insert Radius: ")
	fmt.Scanln(&radius)

	xfloat, _ = strconv.ParseFloat(x, 64)
	yfloat, _ = strconv.ParseFloat(y, 64)

	data.Coordinate = [2]float64{xfloat, yfloat}
	data.Distance, _ = strconv.ParseFloat(distance, 64)
	data.StarDistance, _ = strconv.ParseFloat(stardistance, 64)
	data.Mass, _ = strconv.ParseFloat(mass, 64)
	data.Radius, _ = strconv.ParseFloat(mass, 64)

	stream, _ := json.Marshal(data)
	fmt.Printf("[*] Found a Planet with those specs: %s\n", string(stream))

	sendData(data)
	fmt.Printf("[.] Data sent")
}

func sendData(data interface{}) {
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

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := client.Publish("iot/telescope", 2, false, data)
	if !token.Wait() {
		panic(fmt.Errorf("error sending the message"))
	}
}
