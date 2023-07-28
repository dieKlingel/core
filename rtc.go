package main

import "github.com/pion/webrtc/v3"

type RTC struct {
	Connection *webrtc.PeerConnection
	Tracks     []*webrtc.TrackLocalStaticSample
}
