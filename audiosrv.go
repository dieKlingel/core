package main

import (
	"log"
	"sync"

	"github.com/dieklingel/core/internal/core"
	"github.com/dieklingel/core/internal/io"
)

type AudioService struct {
	storageService core.StorageService

	microphone io.Microphone
	mutex      sync.Mutex
}

const (
	DefaultMicrophonePipeline = "audiotestsrc ! audio/x-raw, format=S16LE, layout=interleaved, rate=48000, channels=1 ! appsink name=rawsink"
)

func NewAudioService(storageService core.StorageService) *AudioService {
	return &AudioService{
		storageService: storageService,
	}
}

func (service *AudioService) MicrophonePipeline() string {
	pipeline := service.storageService.Read().Media.Audio.Src
	if len(pipeline) == 0 {
		return DefaultMicrophonePipeline
	}

	return pipeline
}

func (service *AudioService) NewMicrophoneStream(codec io.AudioCodec) *io.Stream {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	if service.microphone == nil {
		microphone, err := io.NewMicrophone(service.MicrophonePipeline())
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		service.microphone = microphone
	}

	stream, err := service.microphone.NewStream(codec)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return stream
}

func (service *AudioService) ReleaseMicrophoneStream(stream *io.Stream) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	service.microphone.ReleaseStream(stream)
	if !service.microphone.HasStream() {
		service.microphone = nil
	}
}
