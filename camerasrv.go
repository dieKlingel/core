package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"path"
	"time"

	"github.com/dieklingel/core/internal/videosrc"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pion/mediadevices"
)

func RegisterCameraHandler(prefix string, client mqtt.Client) {
	Register(client, path.Join(prefix, "snapshot"), onSnapshot)
}

func onSnapshot(c mqtt.Client, req Request) Response {
	/* stream, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(constraint *mediadevices.MediaTrackConstraints) {},
	})
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("cannot capture frame: %s", err.Error()), 500)
	}*/

	tr, err := videosrc.NewSharedVideoTrack(&mediadevices.CodecSelector{})
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("cannot capture frame: %s", err.Error()), 500)
	}
	videoTrack := tr.(*mediadevices.VideoTrack)
	defer videoTrack.Close()
	// Create a new video reader to get the decoded frames. Release is used
	// to return the buffer to hold frame back to the source so that the buffer
	// can be reused for the next frames.
	time.Sleep(time.Millisecond * 500)

	videoReader := videoTrack.NewReader(false)
	frame, release, _ := videoReader.Read()
	var output bytes.Buffer
	jpeg.Encode(&output, frame, nil)
	release()

	body := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(output.Bytes())
	response := NewResponseFromString(body, 200)
	return response.WithContentType("image/jpeg")
}
