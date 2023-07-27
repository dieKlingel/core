package main

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func RegisterSignHandler(prefix string, client mqtt.Client) {
	Register(client, path.Join(prefix), onSigns)
	Register(client, path.Join(prefix, "ring", "+"), onSignRing)
}

func onSigns(cient mqtt.Client, req Request) Response {
	config, err := NewConfigFromCurrentDirectory()
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not create config: %s.", err.Error()), 500)
	}

	json, err := json.Marshal(config.Gui.Signs)
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not serialize signs: %s", err.Error()), 500)
	}

	return NewResponseFromString(string(json), 200)
}

func onSignRing(client mqtt.Client, req Request) Response {
	pathSegments := strings.Split(req.RequestPath, "/")
	sign := pathSegments[len(pathSegments)-1]

	ExecuteActionsFromPattern(
		"ring",
		map[string]string{
			"SIGN": sign,
		},
	)

	return NewResponseFromString("Ok", 200)
}
