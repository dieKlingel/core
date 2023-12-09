package main

import (
	"github.com/dieklingel/core/audio"
	"github.com/dieklingel/core/camera"
	"github.com/dieklingel/core/internal/core"
	"go.uber.org/fx"
)

func NewFxStorageService(lc fx.Lifecycle) core.StorageService {
	return NewStorageService("core.yaml")
}

func NewFxHttpService(lc fx.Lifecycle, storageService core.StorageService, camera *camera.Camera) *HttpService {
	return NewHttpService(8080, storageService, camera)
}

func NewFxActionService(lc fx.Lifecycle, storageService core.StorageService) core.ActionService {
	return NewActionService(storageService)
}

func NewFxMqttService(lc fx.Lifecycle, storageService core.StorageService, actionService core.ActionService, webrtcService *WebRTCService) *MqttService {
	return NewMqttService(storageService, actionService, webrtcService)
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
