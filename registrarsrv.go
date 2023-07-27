package main

import (
	"encoding/json"
	"log"
	"path"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type HandlerFunction func(mqtt.Client, Request) Response

func Register(client mqtt.Client, channel string, handler HandlerFunction) {
	client.Subscribe(channel, 2, func(c mqtt.Client, m mqtt.Message) {
		var request = Request{}

		if err := json.Unmarshal(m.Payload(), &request); err != nil {
			log.Printf("Could not parse the request: %s. Request will be silently ignored.\r\n", err.Error())
			return
		}

		request.RequestPath = m.Topic()
		var response = handler(client, request)
		if request.IsSocketMessage() {
			return
		}

		json, err := json.Marshal(response)
		if err != nil {
			log.Printf("Could not parse the response: %s. Response will no be sent to the remote.\r\n", err.Error())
			return
		}

		answerChannel, err := request.GetAnswerChannel()
		if err != nil {
			log.Printf("Could not parse the answer channel: %s. Response will no be sent to the remote.\r\n", err.Error())
			return
		}

		// use m.Topic() instead if channel, cause channel could contain wildcards lilke + or # on which we cannot send
		client.Publish(path.Join(m.Topic(), answerChannel), 2, false, string(json))
	})
}
