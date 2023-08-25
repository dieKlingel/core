package main

import (
	"github.com/dieklingel/core/internal/gmedia"
	"github.com/dieklingel/core/internal/io"
	"github.com/pion/webrtc/v3"
)

type RTC struct {
	Connection      *webrtc.PeerConnection
	VideoStream     *io.Stream
	AudioStream     *io.Stream
	RemoteAudioSink *gmedia.RemoteAudioSink
}
