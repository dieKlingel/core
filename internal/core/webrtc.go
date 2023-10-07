package core

import (
	"github.com/pion/webrtc/v3"
)

type Peer struct {
	Id string
}

type PeerHooks struct {
	OnCandidate func(Peer, webrtc.ICECandidateInit)
	OnClose     func(Peer)
}

type WebRTCService interface {
	NewConnection(offer webrtc.SessionDescription, hooks PeerHooks) (*Peer, webrtc.SessionDescription)
	GetConnectionById(id string) *Peer
	AddICECandidate(peer *Peer, candidate webrtc.ICECandidateInit)
	CloseConnection(peer *Peer)
}
