package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"path"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pion/mediadevices"
)

func RegisterCameraHandler(prefix string, client mqtt.Client) {
	Register(client, path.Join(prefix, "snapshot"), onSnapshot)
}

func onSnapshot(c mqtt.Client, req Request) Response {
	stream, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(constraint *mediadevices.MediaTrackConstraints) {},
	})
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("cannot capture frame: %s", err.Error()), 500)
	}

	videoTrack := stream.GetVideoTracks()[0].(*mediadevices.VideoTrack)
	defer videoTrack.Close()

	videoReader := videoTrack.NewReader(false)
	frame, release, _ := videoReader.Read()
	var output bytes.Buffer
	jpeg.Encode(&output, frame, nil)
	release()

	body := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(output.Bytes())
	response := NewResponseFromString(body, 200)
	return response.WithContentType("image/jpeg")
}
