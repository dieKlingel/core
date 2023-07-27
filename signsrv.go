package main

import (
	"path"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func RegisterSignHandler(prefix string, client mqtt.Client) {
	Register(client, path.Join(prefix, "ring", "+"), onSignRing)
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

	return NewResponse("Ok", 200)
}
