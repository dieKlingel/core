package core

import "github.com/dieklingel/core/internal/io"

type AudioService interface {
	MicrophonePipeline() string
	NewMicrophoneStream(codec io.AudioCodec) *io.Stream
	ReleaseMicrophoneStream(stream *io.Stream)
}
