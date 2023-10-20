package main

import (
	"log"
	"sync"

	"github.com/dieklingel/core/internal/io"
)

type CameraService struct {
	storageService *StorageService

	camera io.Camera
	mutex  sync.Mutex
}

const (
	DefaultCameraPipeline = "videotestsrc ! video/x-raw, framerate=30/1, width=1280, height=720 ! appsink name=rawsink"
)

func NewCameraService(storageService *StorageService) *CameraService {
	return &CameraService{
		storageService: storageService,
	}
}

func (service *CameraService) CameraPipeline() string {
	pipeline := service.storageService.Read().Media.Camera.Src
	if len(pipeline) == 0 {
		return DefaultCameraPipeline
	}

	return pipeline
}

func (service *CameraService) NewCameraStream(codec io.CameraCodec) *io.Stream {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	if service.camera == nil {
		camera, err := io.NewCamera(service.CameraPipeline())
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		service.camera = camera
	}

	stream, err := service.camera.NewStream(codec)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return stream
}

func (service *CameraService) ReleaseCameraStream(stream *io.Stream) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	service.camera.ReleaseStream(stream)
	if !service.camera.HasStream() {
		service.camera = nil
	}
}
