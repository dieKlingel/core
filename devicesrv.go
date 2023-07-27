package main

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	clover "github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

func RegisterDeviceHandler(prefix string, client mqtt.Client) {
	Register(client, path.Join(prefix), onDevices)
	Register(client, path.Join(prefix, "create"), onCreateDevice)
	Register(client, path.Join(prefix, "update", "+"), onUpdateDevice)
	Register(client, path.Join(prefix, "delete", "+"), onDeleteDevice)
}

func onDevices(client mqtt.Client, req Request) Response {
	db, err := clover.Open("storage")

	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not open the database: %s", err.Error()), 500)
	}

	defer db.Close()

	if succes, _ := db.HasCollection("devices"); !succes {
		db.CreateCollection("devices")
	}

	docs, err := db.FindAll(query.NewQuery("devices"))
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not fetch the devices: %s", err.Error()), 500)
	}

	devices := make([]Device, 0)
	for _, doc := range docs {
		device := NewDeviceFromMap(doc.ToMap())
		devices = append(devices, device)
	}

	json, err := json.Marshal(devices)
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not serializer the result: %s", err.Error()), 500)
	}

	return NewResponseFromString(string(json), 200)
}

func onCreateDevice(client mqtt.Client, req Request) Response {
	device := Device{}
	if err := json.Unmarshal([]byte(req.Body), &device); err != nil {
		return NewResponseFromString(fmt.Sprintf("could not parse the device: %s", err.Error()), 400)
	}

	if len(device.Token) == 0 {
		return NewResponseFromString("the token cannot be empty", 400)
	}

	db, err := clover.Open("storage")
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not open the database: %s", err.Error()), 500)
	}
	defer db.Close()

	if exists, _ := db.HasCollection("devices"); !exists {
		db.CreateCollection("devices")
	}

	doc, _ := db.FindFirst(query.NewQuery("devices").Where(query.Field("token").Eq(device.Token)))
	if doc == nil {
		doc = document.NewDocument()
	}
	doc.Set("token", device.Token)
	doc.Set("signs", device.Signs)
	doc.SetExpiresAt(time.Now().Add(time.Hour * 24 * 60))

	if err := db.Insert("devices", doc); err != nil {
		return NewResponseFromString(fmt.Sprintf("could not insert device: %s", err.Error()), 400)
	}

	return NewResponseFromString("Ok", 201)
}

func onUpdateDevice(client mqtt.Client, req Request) Response {
	pathSegments := strings.Split(req.RequestPath, "/")
	token := pathSegments[len(pathSegments)-1]

	device := Device{}
	if err := json.Unmarshal([]byte(req.Body), &device); err != nil {
		return NewResponseFromString(fmt.Sprintf("could not parse the device: %s", err.Error()), 400)
	}

	if len(device.Token) == 0 {
		return NewResponseFromString("the token cannot be empty", 400)
	}

	db, err := clover.Open("storage")
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not open the database: %s", err.Error()), 500)
	}
	defer db.Close()

	values := make(map[string]interface{}, 0)
	values["token"] = device.Token
	values["signs"] = device.Signs

	if err := db.Update(query.NewQuery("devices").Where(query.Field("token").Eq(token)), values); err != nil {
		return NewResponseFromString(fmt.Sprintf("could not update the device: %s", err.Error()), 400)
	}

	return NewResponseFromString("Ok", 200)
}

func onDeleteDevice(client mqtt.Client, req Request) Response {
	pathSegments := strings.Split(req.RequestPath, "/")
	token := pathSegments[len(pathSegments)-1]

	db, err := clover.Open("storage")
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not open the database: %s", err.Error()), 500)
	}
	defer db.Close()

	if err := db.Delete(query.NewQuery("devices").Where(query.Field("token").Eq(token))); err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not delete devices: %s", err.Error()), 500)
	}
	return NewResponseFromString("Ok", 204)
}
