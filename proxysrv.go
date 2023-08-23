package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Access-Control-Allow-Methods", "*")
	writer.Header().Add("Access-Control-Allow-Headers", "*")
	if req.Method == "OPTIONS" {
		writer.WriteHeader(http.StatusOK)
		return
	}

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

	client := mqtt.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Proxy: Could not connect tho the broker. Message: %s\r\n", token.Error().Error())
		writer.WriteHeader(400)
		writer.Write([]byte(token.Error().Error()))
		return
	}

	subTopic := path.Join("./", req.URL.Path, answerChannel.String())
	client.Subscribe(subTopic, 2, func(c mqtt.Client, m mqtt.Message) {
		response := NewEmptyResponse()
		json.Unmarshal(m.Payload(), &response)
		result <- response
	})

	pubTopic := path.Join("./", req.URL.Path)
	client.Publish(pubTopic, 2, false, httpRequestToMqttRequestPayload(req, answerChannel.String()))

	select {
	case res := <-result:
		for header, value := range res.Headers {
			fmt.Print(header + ":" + value)
			writer.Header().Set(header, value)
		}
		writer.WriteHeader(res.StatusCode)
		writer.Write([]byte(res.Body))
	case <-time.After(20 * time.Second):
		writer.WriteHeader(http.StatusNotFound)
	}

	defer client.Disconnect(0)
}

func httpRequestToMqttRequestPayload(req *http.Request, answerChannel string) string {
	request := Request{}

	bytes := make([]byte, req.ContentLength)
	req.Body.Read(bytes)
	request.Body = string(bytes)

	request.Method = req.Method
	request.Headers = make(map[string]string)
	request.Headers["mqtt_answer_channel"] = answerChannel

	result, _ := json.Marshal(request)
	return string(result)
}
