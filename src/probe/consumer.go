package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nuclio/nuclio-sdk-go"
)

type ProbeData struct {
	Name        string
	Humidity    float32
	Temperature float32
	Wind        float32
}

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	logger, err := InitLogger()
	if err != nil {
		context.Logger.Error("Error: %s", err)
		panic(err)
	}
	defer logger.Close()

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

		var data ProbeData
		json.Unmarshal(body, &data)

		err := persistence.PersistProbeData(data)
		if err != nil {
			logger.Error(err.Error())
		} else {
			logger.Info(fmt.Sprintf("Planet `%s` data updated.", data.Name))
		}

		return nil, nil
	}

	return nuclio.Response{
		StatusCode:  http.StatusForbidden,
		ContentType: "text/plain",
		Body:        []byte("This is not a Web API"),
	}, nil
}
