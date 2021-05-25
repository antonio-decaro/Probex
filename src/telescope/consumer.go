package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nuclio/nuclio-sdk-go"
)

const USERNAME_ENV string = "MQTT_USERNAME"
const PASSWORD_ENV string = "MQTT_PASSWORD"
const IP_ENV string = "MQTT_BROKER_IP"
const PORT_ENV string = "PORT_ENV"
const LOG_QUEUE_NAME string = "iot/logs"

type TelescopeData struct {
	Name       string
	Coordinate [2]float64
	Distance   float32
	Mass       int32
	Radius     float32
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
		logger.Log(err.Error())
		context.Logger.Error("Error: %s", err)
		panic(err)
	}

	// if we got the event from rabbit
	if event.GetTriggerInfo().GetClass() == "async" && event.GetTriggerInfo().GetKind() == "mqtt" {
		body := event.GetBody()
		logger.Log("Body content: " + string(body))

		var data SpaceProbeData
		json.Unmarshal(body, &data)

		ch := make(chan error)

		go func(ch chan<- error) {
			var err error
			if classificator.ClassifyData(data) {
				logger.Log(fmt.Sprintf("Planet %+v is potentially habitable!", data))
				err = sendPlanetProbe(data)
			} else {
				logger.Log(fmt.Sprintf("Planet %+v is not habitable", data))
			}
			ch <- err
		}(ch)

		PersistProbeData(data)

		// waiting for classifying operation to complete
		<-ch

		return nil, nil
	}

	return nuclio.Response{
		StatusCode:  http.StatusForbidden,
		ContentType: "text/plain",
		Body:        []byte("This is not a Web API"),
	}, nil
}

func sendPlanetProbe(data SpaceProbeData) error {
	return nil
}
