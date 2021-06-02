package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nuclio/nuclio-sdk-go"
)

type TelescopeData struct {
	Name         string
	Coordinate   [2]float64
	Distance     float32
	StarDistance float32 // astronomic units
	StarType     string
	Mass         float32
	Radius       float32
}

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	logger, err := InitLogger()
	if err != nil {
		context.Logger.Error("Error: %s", err)
		panic(err)
	}
	defer logger.Close()
	classificator, err := InitClassificator()
	if err != nil {
		logger.Error(err.Error())
		context.Logger.Error("Error: %s", err)
		panic(err)
	}
	persistence, err := InitPersistence()
	if err != nil {
		logger.Error(err.Error())
		context.Logger.Error("Error: %s", err)
		panic(err)
	}

	// if we got the event from rabbit
	if event.GetTriggerInfo().GetClass() == "async" && event.GetTriggerInfo().GetKind() == "mqtt" {
		body := event.GetBody()
		logger.Debug("Body content: " + string(body))

		var data TelescopeData
		json.Unmarshal(body, &data)

		ch := make(chan error)

		go func(ch chan<- error) {
			var err error
			if classificator.ClassifyData(data) {
				logger.Info(fmt.Sprintf("Planet %+v is potentially habitable! Sending Probe...", data))
				err = sendPlanetProbe(data)
			} else {
				logger.Info(fmt.Sprintf("Planet %+v is not habitable.", data))
			}
			ch <- err
		}(ch)

		err := persistence.PersistTelescopeData(data)
		if err != nil {
			logger.Error(err.Error())
		} else {
			logger.Info("Planet successfully inserted in the database.")
		}

		// waiting for classifying operation to complete
		if err := <-ch; err != nil {
			logger.Error(err.Error())
		}

		return nil, nil
	}

	return nuclio.Response{
		StatusCode:  http.StatusForbidden,
		ContentType: "text/plain",
		Body:        []byte("This is not a Web API"),
	}, nil
}

func sendPlanetProbe(data TelescopeData) error {
	client := newMQTTClient("mqtt_probe_sender")
	defer client.Disconnect(250)

	msg := map[string]interface{}{
		"Name":       data.Name,
		"Coordinate": data.Coordinate,
		"Distance":   data.Distance,
	}

	send, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	token := client.Publish(PROBE_TOPIC_NAME, 2, false, send)
	if !token.Wait() {
		return fmt.Errorf("error sending the message")
	}

	return nil
}

func newMQTTClient(id string) mqtt.Client {
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
