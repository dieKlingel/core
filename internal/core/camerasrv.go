package core

import "github.com/dieklingel/core/internal/io"

type CameraService interface {
	CameraPipeline() string
	SetCameraPipeline(pipeline string)
	NewCameraStream(codec io.CameraCodec) *io.Stream
	ReleaseCameraStream(stream *io.Stream)
}
