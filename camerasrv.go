package main

import (
	"encoding/base64"
	"path"

	"github.com/dieklingel/core/internal/io"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func RegisterCameraHandler(prefix string, client mqtt.Client) {
	Register(client, path.Join(prefix, "snapshot"), onSnapshot)
}

func onSnapshot(c mqtt.Client, req Request) Response {
	if camera == nil {
		return NewResponseFromString("the camera was not created succesfully, chech your logs for more information", 500)
	}

	stream, err := io.NewStream("appsrc name=src ! videoconvert ! pngenc ! appsink sync=false name=sink")
	if err != nil {
		return NewResponseFromString(err.Error(), 500)
	}

	camera.AddStream(stream)

	select {
	case <-stream.Finished:
		return NewResponseFromString("could not receive a frame from camera", 500)
	case frame := <-stream.Frame:
		camera.RemoveStream(stream)
		return NewResponseFromString("data:image/png;base64,"+base64.StdEncoding.EncodeToString(frame.GetBuffer().Bytes()), 200)
	}
}
