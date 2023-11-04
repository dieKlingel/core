package main

import (
	"fmt"
	"time"

	"github.com/dieklingel/core/internal/core"
	"github.com/dieklingel/core/internal/io"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

type Peer struct {
	connection  *webrtc.PeerConnection
	videostream *io.Stream
	audiostream *io.Stream
}

type WebRTCService struct {
	CameraService *CameraService
	audioService  core.AudioService

	connections map[string]*Peer
}

func NewWebRTCService(camerasrv *CameraService, audioService core.AudioService) *WebRTCService {
	return &WebRTCService{
		CameraService: camerasrv,
		audioService:  audioService,
	}
}

func (service *WebRTCService) NewConnection(offer webrtc.SessionDescription, hooks core.PeerHooks) (*core.Peer, webrtc.SessionDescription) {
	peerConnection, err := webrtc.NewPeerConnection(
		webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
			SDPSemantics: webrtc.SDPSemanticsUnifiedPlan,
		},
	)
	if err != nil && hooks.OnClose != nil {
		hooks.OnClose(core.Peer{})
	}
	peer := core.Peer{
		Id: uuid.New().String(),
	}
	service.connections[peer.Id] = &Peer{
		connection: peerConnection,
	}

	// videotrack
	videotrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, fmt.Sprintf("video-%s", uuid.New().String()), "pion-video")
	if err != nil {
		panic(err)
	}
	_, err = peerConnection.AddTrack(videotrack)
	if err != nil {
		println("could not add track" + err.Error())
	}
	videostream := service.CameraService.NewCameraStream(io.X264CameraCodec)

	go func() {
		for {
			select {
			case sample := <-videostream.Frame:
				videotrack.WriteSample(media.Sample{
					Data:     sample.GetBuffer().Bytes(),
					Duration: 1 * time.Millisecond, // use 1ms, because duration is incorrect when used with libcamerasrc, which is our preffered way
				})
			case <-videostream.Finished:
				return
			}
		}
	}()

	// audiotrack
	audiotrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, fmt.Sprintf("audio-%s", uuid.New().String()), "pion-audio")
	if err != nil {
		panic(err)
	}
	_, err = peerConnection.AddTrack(audiotrack)
	if err != nil {
		println("could not add track" + err.Error())
	}
	audiostream := service.audioService.NewMicrophoneStream(io.OpusAudioCodec)

	go func() {
		for {
			select {
			case sample := <-audiostream.Frame:
				audiotrack.WriteSample(media.Sample{
					Data:     sample.GetBuffer().Bytes(),
					Duration: sample.GetBuffer().Duration(),
				})
			case <-audiostream.Finished:
				return
			}
		}
	}()

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if hooks.OnCandidate != nil {
			hooks.OnCandidate(peer, i.ToJSON())
		}
	})

	// TODO: on track

	peerConnection.SetRemoteDescription(offer)
	answer, _ := peerConnection.CreateAnswer(&webrtc.AnswerOptions{})
	peerConnection.SetLocalDescription(answer)

	return &peer, answer
}

func (service *WebRTCService) GetConnectionById(id string) *core.Peer {
	if _, exists := service.connections[id]; exists {
		return &core.Peer{
			Id: id,
		}
	}

	return nil
}

func (service *WebRTCService) AddICECandidate(peer *core.Peer, candidate webrtc.ICECandidateInit) {
	if p, exists := service.connections[peer.Id]; exists {
		p.connection.AddICECandidate(candidate)
	}
}

func (service *WebRTCService) CloseConnection(peer *core.Peer) {
	if p, exists := service.connections[peer.Id]; exists {
		p.connection.Close()
		if p.videostream != nil {
			service.CameraService.ReleaseCameraStream(p.videostream)
		}
		if p.audiostream != nil {
			service.audioService.ReleaseMicrophoneStream(p.audiostream)
		}
	}
}
