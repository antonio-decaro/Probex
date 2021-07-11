package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	table "github.com/jedib0t/go-pretty/v6/table"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

const USERNAME_ENV string = "MQTT_USERNAME"
const PASSWORD_ENV string = "MQTT_PASSWORD"
const IP_ENV string = "MQTT_BROKER_IP"
const PORT_ENV string = "PORT_ENV"
const LOG_QUEUE_NAME string = "iot/monitor"
const FNAME string = "./monitor.dat"

var (
	planetinfo []map[string]interface{}
	headers    = []string{
		"Name",
		"Coordinate",
		"Distance",
		"StarDistance",
		"StarType",
		"Mass",
		"Radius",
		"Humidity",
		"Temperature",
		"Wind",
		"ProbeId",
	}
	writer table.Writer
)

func main() {
	writer = table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	header := table.Row{"#"}
	for _, h := range headers {
		header = append(header, h)
	}

	writer.AppendHeader(header)
	writer.SetStyle(table.StyleColoredBright)

	loadData()

	// to save the file on exit
	ctrlc := make(chan os.Signal)
	signal.Notify(ctrlc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctrlc
		saveData()
		os.Exit(0)
	}()

	fmt.Println("[*] Starting monitor...")
	printData()
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
		// fmt.Printf("%+v\n", data)
		found := false

		for _, el := range planetinfo {
			if el["Name"] == data["Name"] {
				found = true
				for key, val := range data {
					el[key] = val
				}
				break
			}
		}

		if !found {
			planetinfo = append(planetinfo, data)
		}

		printData()
		saveData()

	} else {
		fmt.Println(err)
	}
}

func printData() {

	fmt.Print("\033[H\033[2J")
	writer.ResetRows()

	for i, el := range planetinfo {
		row := table.Row{i}
		for _, h := range headers {
			if val, ok := el[h]; ok {
				row = append(row, val)
			} else {
				row = append(row, "")
			}
		}
		writer.AppendRow(row)
	}

	writer.Render()
}

func loadData() {
	data, err := ioutil.ReadFile(FNAME)
	if err != nil {
		fmt.Println(err)
		return
	}

	json.Unmarshal(data, &planetinfo)
}

func saveData() {
	data, err := json.Marshal(planetinfo)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(FNAME, data, 0777)
	if err != nil {
		panic(err)
	}
}
