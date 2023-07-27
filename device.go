package main

import "fmt"

type Device struct {
	Token string   `json:"token" clover:"token"`
	Signs []string `json:"signs" clover:"signs"`
}

func NewDevice(token string, signs []string) Device {
	return Device{
		Token: token,
		Signs: signs,
	}
}

func NewDeviceFromMap(values map[string]interface{}) Device {
	token := fmt.Sprintf("%s", values["token"])
	signs := make([]string, 0)

	for _, sign := range values["signs"].([]interface{}) {
		signs = append(signs, fmt.Sprintf("%s", sign))
	}

	return Device{
		Token: token,
		Signs: signs,
	}
}
