package main

import (
	"fmt"
	"log"

	"github.com/dieklingel/core/audio"
	"github.com/dieklingel/core/camera"
	"github.com/dieklingel/core/internal/core"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

type Peer struct {
	connection  *webrtc.PeerConnection
	videostream *camera.CameraStream
	audiostream *audio.AudioStream
}

type WebRTCService struct {
	camera     *camera.Camera
	audioInput *audio.Input

	connections map[string]*Peer
}

func NewWebRTCService(camera *camera.Camera, audioInput *audio.Input) *WebRTCService {
	return &WebRTCService{
		camera:     camera,
		audioInput: audioInput,

		connections: make(map[string]*Peer),
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
	videostream := service.camera.Tee(camera.X264CameraCodec)
	service.connections[peer.Id].videostream = videostream

	go func() {
		for {
			select {
			case sample := <-videostream.Frame():
				videotrack.WriteSample(media.Sample{
					Data:     sample.GetBuffer().Bytes(),
					Duration: *sample.GetBuffer().Duration().AsDuration(), // use 1ms, because duration is incorrect when used with libcamerasrc, which is our preffered way
				})
				// TODO timeout or emit finish
			case <-videostream.Finished():
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
	audiostream := service.audioInput.Tee(audio.OpusAudioCodec) //service.audioService.NewMicrophoneStream(io.OpusAudioCodec)
	service.connections[peer.Id].audiostream = audiostream

	go func() {
		for {
			select {
			case sample := <-audiostream.Frame():
				audiotrack.WriteSample(media.Sample{
					Data:     sample.GetBuffer().Bytes(),
					Duration: *sample.GetBuffer().Duration().AsDuration(),
				})
			case <-audiostream.Finished():
				return
			}
		}
	}()

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			log.Println("ice candidate was nil")
			return
		}
		if hooks.OnCandidate != nil {
			hooks.OnCandidate(peer, i.ToJSON())
		}
	})

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		if state == webrtc.ICEConnectionStateDisconnected {
			if hooks.OnClose != nil {
				hooks.OnClose(peer)
			}
			service.CloseConnection(&peer)
			return
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
			p.videostream.FinishAndClose()
		}
		if p.audiostream != nil {
			p.audiostream.FinishAndClose()
		}
	}
	delete(service.connections, peer.Id)
}
