package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/dieklingel/core/internal/gmedia"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
)

var connections map[string]RTC = make(map[string]RTC)

func RegisterRtcHandler(prefix string, client mqtt.Client) {
	Register(client, path.Join(prefix, "connections"), onGetConnections)
	Register(client, path.Join(prefix, "connections", "create", "+"), onCreateConnection)
	Register(client, path.Join(prefix, "connections", "close", "+"), onCloseConnection)
	Register(client, path.Join(prefix, "connections", "candidate", "+"), onAddCandidate)
}

func onGetConnections(c mqtt.Client, req Request) Response {
	json, err := json.Marshal(connections)
	if err != nil {
		return NewResponseFromString(fmt.Sprintf("Could not encode: %s.", err), 500)
	}

	return NewResponseFromString(string(json), 200)
}

func onCreateConnection(client mqtt.Client, req Request) Response {
	pathSegments := strings.Split(req.RequestPath, "/")
	id := pathSegments[len(pathSegments)-1]

	if _, exists := connections[id]; exists {
		return NewResponseFromString(fmt.Sprintf("Cannot create a connection with id '%s' because a connection with this id already exists.", id), 409)
	}

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

	rtc := RTC{
		Connection: peerConnection,
		Tracks:     make([]*webrtc.TrackLocalStaticSample, 0),
	}

	if err != nil {
		peerConnection.Close()
		return NewResponseFromString(fmt.Sprintf("Cannot create a connection: %s", err.Error()), 500)
	}

	firstVideoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, fmt.Sprintf("video-%s", uuid.New().String()), "pion2")
	if err != nil {
		panic(err)
	}

	rtc.Tracks = append(rtc.Tracks, firstVideoTrack)
	gmedia.AddVideoTrack(firstVideoTrack)

	_, err = peerConnection.AddTrack(firstVideoTrack)
	if err != nil {
		panic(err)
	}

	firstAudioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, fmt.Sprintf("audio-%s", uuid.New().String()), "pion3")
	if err != nil {
		panic(err)
	}

	rtc.Tracks = append(rtc.Tracks, firstAudioTrack)
	gmedia.AddAudioTrack(firstAudioTrack)

	_, err = peerConnection.AddTrack(firstAudioTrack)
	if err != nil {
		panic(err)
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			log.Println("OnICECandidate Callback was called with nil reference and will be silently ignored.")
			return // dont know why the callback is called with nil, but it is the case
		}
		candidate, _ := json.Marshal(c.ToJSON())
		request := NewSocketRequest(string(candidate))
		request.Method = "CONNECT"
		json, _ := json.Marshal(request)

		client.Publish(path.Join(id, "connection", "candidate"), 2, false, string(json))
	})

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		fmt.Printf("Track has started, of type %d: %s \n", track.PayloadType(), track.Codec().MimeType)
		// TODO: play stream, if audio
	})

	var offer *webrtc.SessionDescription = &webrtc.SessionDescription{}
	if err := json.Unmarshal([]byte(req.Body), offer); err != nil {
		panic(err.Error())
	}
	if err := peerConnection.SetRemoteDescription(*offer); err != nil {
		panic(err.Error())
	}

	answer, err := peerConnection.CreateAnswer(&webrtc.AnswerOptions{})
	peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err.Error())
	}

	json, err := json.Marshal(answer)
	if err != nil {
		panic(err.Error())
	}

	connections[id] = rtc

	return NewResponseFromString(string(json), 201)
}

func onCloseConnection(client mqtt.Client, req Request) Response {
	pathSegments := strings.Split(req.RequestPath, "/")
	id := pathSegments[len(pathSegments)-1]

	if rtc, exists := connections[id]; exists {
		rtc.Connection.Close()
		for _, track := range rtc.Tracks {
			gmedia.RemoveAudioTrack(track)
			gmedia.RemoveVideoTrack(track)

		}
	}

	return NewResponseFromString("", 200)
}

func onAddCandidate(client mqtt.Client, req Request) Response {
	pathSegments := strings.Split(req.RequestPath, "/")
	id := pathSegments[len(pathSegments)-1]

	candidate := &webrtc.ICECandidateInit{}
	if err := json.Unmarshal([]byte(req.Body), candidate); err != nil {
		log.Printf("could not parse candidate: %s", err.Error())
		return NewResponseFromString(fmt.Sprintf("the canidate could not be parsed: %s", err.Error()), 400)
	}

	if rtc, exists := connections[id]; exists {
		if err := rtc.Connection.AddICECandidate(*candidate); err != nil {
			log.Printf("could not add candidate: %s", err.Error())
			return NewResponseFromString(fmt.Sprintf("could not add the candidate: %s", err.Error()), 500)
		}
		return NewResponseFromString("ok", 200)
	}

	return NewResponseFromString("the requested resource does not exist", 404)
}