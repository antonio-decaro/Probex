package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	persistence, err := InitPersistence(logger)
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
				logger.Info(fmt.Sprintf("Planet %+v is potentially habitable!", data))
				err = sendPlanetProbe(data)
			} else {
				logger.Info(fmt.Sprintf("Planet %+v is not habitable", data))
			}
			ch <- err
		}(ch)

		persistence.PersistTelescopeData(data)

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
	return nil
}
