package main

import (
	"net/http"

	"github.com/nuclio/nuclio-sdk-go"
)

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	logger, err := InitLogger()
	if err != nil {
		context.Logger.Error("Error: %s", err)
		panic(err)
	}
	defer logger.Close()

	// if we got the event from rabbit
	if event.GetTriggerInfo().GetClass() == "async" && event.GetTriggerInfo().GetKind() == "mqtt" {
		body := event.GetBody()
		logger.Debug("Body content: " + string(body))

		// TODO data with probe

		return nil, nil
	}

	return nuclio.Response{
		StatusCode:  http.StatusForbidden,
		ContentType: "text/plain",
		Body:        []byte("This is not a Web API"),
	}, nil
}
