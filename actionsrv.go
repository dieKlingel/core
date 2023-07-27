package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path"
	"regexp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func RegisterActionHandler(prefix string, client mqtt.Client) {
	Register(client, path.Join(prefix), onActions)
	Register(client, path.Join(prefix, "execute"), onExecuteActions)
}

func onActions(client mqtt.Client, req Request) Response {
	config, err := NewConfigFromCurrentDirectory()
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not read config: %s", err.Error()), 500)
	}

	json, err := json.Marshal(config.Actions)
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not serialize actions: %s", err.Error()), 500)
	}

	return NewResponseFromString(string(json), 200)
}

func onExecuteActions(client mqtt.Client, req Request) Response {
	payload := make(map[string]interface{})
	json.Unmarshal([]byte(req.Body), &payload)

	pattern, ok := payload["pattern"].(string)
	if !ok {
		return NewResponseFromString("the pattern has to be of type string", 400)
	}

	env, ok := payload["environment"].(map[string]interface{})
	if !ok {
		return NewResponseFromString("the evironment has to be of type {string: string}", 400)
	}
	environment := make(map[string]string)
	for key, value := range env {
		environment[key] = fmt.Sprintf("%s", value)
	}

	actions := ExecuteActionsFromPattern(pattern, environment)
	json, _ := json.Marshal(actions)
	return NewResponseFromString(string(json), 200)
}

func ExecuteActionsFromPattern(pattern string, environment map[string]string) []Action {
	config, err := NewConfigFromCurrentDirectory()
	if err != nil {
		log.Printf("could not execute actions: %s", err.Error())
		return make([]Action, 0)
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("error while compiling regex: %s", err.Error())
		return make([]Action, 0)
	}

	actions := make([]Action, 0)
	for _, action := range config.Actions {
		match := regex.MatchString(action.Trigger)
		if match {
			actions = append(actions, action)
			if err := exec.Command("bash", "-c", action.Lane).Run(); err != nil {
				log.Printf("error while running action: %s", err.Error())
			}
		}
	}

	return actions
}
