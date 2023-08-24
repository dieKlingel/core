package main

import (
	"github.com/dieklingel/core/internal/gmedia"
	"github.com/dieklingel/core/internal/video"
	"github.com/pion/webrtc/v3"
)

type RTC struct {
	Connection      *webrtc.PeerConnection
	VideoStream     *video.Stream
	AudioTracks     []*webrtc.TrackLocalStaticSample
	RemoteAudioSink *gmedia.RemoteAudioSink
}
