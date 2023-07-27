package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func RunProxy(port int) {
	go func() {
		http.HandleFunc("/", proxy)
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}

func proxy(writer http.ResponseWriter, req *http.Request) {
	id := uuid.New()
	config, err := NewConfigFromCurrentDirectory()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	url, err := url.Parse(config.Mqtt.Uri)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	answerChannel := uuid.New()
	result := make(chan Response, 1)

	var options = mqtt.NewClientOptions()
	options.AddBroker(url.String())
	options.SetClientID(id.String())
	options.SetUsername(strings.Join(req.Header["Username"], ""))
	options.SetPassword(strings.Join(req.Header["Password"], ""))
	options.SetAutoReconnect(false)
	options.OnConnect = func(c mqtt.Client) {
		subTopic := path.Join("./", req.URL.Path, answerChannel.String())
		c.Subscribe(subTopic, 2, func(c mqtt.Client, m mqtt.Message) {
			response := NewEmptyResponse()
			json.Unmarshal(m.Payload(), &response)
			result <- response
		})

		pubTopic := path.Join("./", req.URL.Path)
		c.Publish(pubTopic, 2, false, httpRequestToMqttRequestPayload(*req, answerChannel.String()))
	}

	client := mqtt.NewClient(options)
	client.Connect()

	select {
	case res := <-result:
		writer.WriteHeader(res.StatusCode)
		writer.Write([]byte(res.Body))
	case <-time.After(30 * time.Second):
		writer.WriteHeader(http.StatusNotFound)
	}
	client.Disconnect(0)
}

func httpRequestToMqttRequestPayload(req http.Request, answerChannel string) string {
	request := Request{}
	req.Body.Read([]byte(request.Body))
	request.Method = req.Method
	request.Headers = make(map[string]string)
	request.Headers["mqtt_answer_channel"] = answerChannel

	result, _ := json.Marshal(request)
	return string(result)
}
