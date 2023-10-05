package core

import "github.com/dieklingel/core/internal/io"

type CameraService interface {
	NewCameraStream() *io.Stream
	ReleaseCameraStream(stream *io.Stream)
}
