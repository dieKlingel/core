package main

import (
	"github.com/dieklingel/core/audio"
	"github.com/dieklingel/core/camera"
)

func NewFxHttpService(camera *camera.Camera) *HttpService {
	return NewHttpService(8080, camera)
}

func NewFxCamera() *camera.Camera {
	camera, err := camera.New("videotestsrc ! video/x-raw, framerate=30/1, width=1280, height=720 ! appsink name=rawsink")
	if err != nil {
		panic(err)
	}

	return camera
}

func NewFxAudioInput() *audio.Input {
	input, err := audio.NewInput("audiotestsrc ! audio/x-raw, format=S16LE, layout=interleaved, rate=48000, channels=1 ! appsink name=rawsink")
	if err != nil {
		panic(err)
	}

	return input
}
