package main

import (
	"github.com/dieklingel/core/internal/gmedia"
	"github.com/pion/webrtc/v3"
)

type RTC struct {
	Connection      *webrtc.PeerConnection
	VideoTracks     []*webrtc.TrackLocalStaticSample
	AudioTracks     []*webrtc.TrackLocalStaticSample
	RemoteAudioSink *gmedia.RemoteAudioSink
}
