package core

import "github.com/dieklingel/core/internal/io"

type MicrophoneService interface {
	NewMicrophoneStream() *io.Stream
	ReleaseMicrophoneStream(stream *io.Stream)
}
