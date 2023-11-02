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

func NewFxHttpService(lc fx.Lifecycle, storageService core.StorageService, cameraService *CameraService) *HttpService {
	return NewHttpService(8080, storageService, cameraService)
}

func NewFxActionService(lc fx.Lifecycle, storageService core.StorageService) *ActionService {
	return NewActionService(storageService)
}

func NewFxWebRTCService(lc fx.Lifecycle, cameraService *CameraService) *WebRTCService {
	return NewWebRTCService(cameraService)
}

func NewFxMqttService(lc fx.Lifecycle, storageService core.StorageService, actionService *ActionService, webrtcService *WebRTCService, appService core.AppService) *MqttService {
	return NewMqttService(storageService, actionService, webrtcService, appService)
}

func NewFxPluginService(lc fx.Lifecycle, storageService *StorageService) *PluginService {
	return NewPluginService(storageService)
}

func NewFxAppService(lc fx.Lifecycle, storageService core.StorageService) core.AppService {
	return NewAppService(storageService)
}
