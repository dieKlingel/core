package main

import (
	"github.com/dieklingel/core/internal/core"
	"go.uber.org/fx"
)

func NewFxStorageService(lc fx.Lifecycle) core.StorageService {
	return NewStorageService("core.yaml")
}

func NewFxCameraService(lc fx.Lifecycle, storageService core.StorageService) *CameraService {
	return NewCameraService(storageService)
}

func NewFxAudioService(lc fx.Lifecycle, storageService core.StorageService) core.AudioService {
	return NewAudioService(storageService)
}

func NewFxHttpService(lc fx.Lifecycle, storageService core.StorageService, cameraService *CameraService) *HttpService {
	return NewHttpService(8080, storageService, cameraService)
}

func NewFxActionService(lc fx.Lifecycle, storageService core.StorageService) core.ActionService {
	return NewActionService(storageService)
}

func NewFxWebRTCService(lc fx.Lifecycle, cameraService *CameraService, audioService core.AudioService) *WebRTCService {
	return NewWebRTCService(cameraService, audioService)
}

func NewFxMqttService(lc fx.Lifecycle, storageService core.StorageService, actionService core.ActionService, webrtcService *WebRTCService) *MqttService {
	return NewMqttService(storageService, actionService, webrtcService)
}

func NewFxPluginService(lc fx.Lifecycle, storageService *StorageService) *PluginService {
	return NewPluginService(storageService)
}
