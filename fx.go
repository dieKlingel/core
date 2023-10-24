package main

import (
	"go.uber.org/fx"
)

func NewFxStorageService(lc fx.Lifecycle) *StorageService {
	return NewStorageService("core.yaml")
}

func NewFxCameraService(lc fx.Lifecycle, storageService *StorageService) *CameraService {
	return NewCameraService(storageService)
}

func NewFxHttpService(lc fx.Lifecycle, storageService *StorageService, cameraService *CameraService) *HttpService {
	return NewHttpService(8080, storageService, cameraService)
}

func NewFxActionService(lc fx.Lifecycle, storageService *StorageService) *ActionService {
	return NewActionService(storageService)
}

func NewFxWebRTCService(lc fx.Lifecycle, cameraService *CameraService) *WebRTCService {
	return NewWebRTCService(cameraService)
}

func NewFxMqttService(lc fx.Lifecycle, storageService *StorageService, actionService *ActionService, webrtcService *WebRTCService) *MqttService {
	return NewMqttService(storageService, actionService, webrtcService)
}
