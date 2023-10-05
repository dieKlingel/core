package main

import (
	"log"
	"net/url"
	"path"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func RunApi(url url.URL, username string, password string) {
	var id = uuid.New()
	log.Printf("Start connection with id '%s'.\r\n", id.String())
	prefix := path.Join("./", url.Path)

	var options = mqtt.NewClientOptions()
	options.AddBroker(url.String())
	options.SetClientID(id.String())
	options.SetUsername(username)
	options.SetPassword(password)
	options.SetAutoReconnect(true)
	options.SetKeepAlive(10 * time.Second)
	options.OnConnect = onConnect
	options.OnConnectionLost = onConnectionLost

	var client = mqtt.NewClient(options)

	for !client.IsConnected() {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("Could not connect tho the broker. Message: %s\r\n", token.Error().Error())
			log.Printf("Retry connecting in 10 seconds.\r\n")
			time.Sleep(10 * time.Second)
		}
	}

	RegisterCameraHandler(prefix+"/camera", client)
	RegisterRtcHandler(prefix+"/rtc", client)
	RegisterDeviceHandler(prefix+"/devices", client)
}

func onConnect(client mqtt.Client) {
	log.Printf("Connected to the broker.\r\n")
}

func onConnectionLost(client mqtt.Client, err error) {
	log.Printf("Connection lost.\r\n")
}
