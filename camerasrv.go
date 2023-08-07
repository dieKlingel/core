package main

import (
	"path"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func RegisterCameraHandler(prefix string, client mqtt.Client) {
	Register(client, path.Join(prefix, "snapshot"), onSnapshot)
}

func onSnapshot(c mqtt.Client, req Request) Response {
	return NewResponseFromString("not implemented", 501)
}
