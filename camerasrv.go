package main

import (
	"bytes"
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

	// Since track can represent audio as well, we need to cast it to
	// *mediadevices.VideoTrack to get video specific functionalities
	track := stream.GetVideoTracks()[0]
	videoTrack := track.(*mediadevices.VideoTrack)
	defer videoTrack.Close()
	// Create a new video reader to get the decoded frames. Release is used
	// to return the buffer to hold frame back to the source so that the buffer
	// can be reused for the next frames.
	videoReader := videoTrack.NewReader(false)
	frame, release, _ := videoReader.Read()
	var output bytes.Buffer
	jpeg.Encode(&output, frame, nil)
	release()

	response := NewResponse(output.Bytes(), 200)
	return response.WithContentType("image/jpeg")
}
