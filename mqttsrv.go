package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/dieklingel/core/config"
	"github.com/dieklingel/core/internal/core"
	"github.com/dieklingel/core/internal/mqtt"
	"github.com/pion/webrtc/v3"
)

type MqttService struct {
	config        *config.Environment
	actionService *ActionService
	webRTCService *WebRTCService

	client *mqtt.Client
}

func NewMqttService(config *config.Environment, actionsrv *ActionService, webrtcsrc *WebRTCService) *MqttService {
	return &MqttService{
		config:        config,
		actionService: actionsrv,
		webRTCService: webrtcsrc,
	}
}

func (service *MqttService) Run() {
	if service.client != nil {
		service.client.Disconnect()
	}

	url := service.config.Mqtt.Uri
	username := service.config.Mqtt.Username
	password := service.config.Mqtt.Password

	service.client = mqtt.NewClient()
	service.client.SetAutoReconnect(true)
	service.client.SetBroker(url)
	service.client.SetUsername(username)
	service.client.SetPassword(password)
	go func() {
		service.client.Connect()

		for !service.client.IsConnected() {
			log.Printf("could not connect to %s; retry in 10 src", url)
			time.Sleep(10 * time.Second)
			service.client.Connect()
		}

		service.buildListeners(service.client, username)
	}()
}

func (service *MqttService) buildListeners(client *mqtt.Client, prefix string) {
	type Headers struct {
		SenderDeviceId  string `json:"senderDeviceId"`
		SenderSessionId string `json:"senderSessionId"`
		SessionId       string `json:"sessionId"`
	}

	type Body struct {
		SessionDescription webrtc.SessionDescription `json:"sessionDescription"`
		IceCandidate       webrtc.ICECandidateInit   `json:"iceCandidate"`
	}

	type ConnectionDescriptionMessage struct {
		Headers Headers `json:"header"`
		Body    Body    `json:"body"`
	}

	type ConnectionCandidateMessage struct {
		Headers Headers `json:"header"`
		Body    Body    `json:"body"`
	}

	type ConnectionCloseMessage struct {
		Headers Headers `json:"header"`
	}

	client.Subscribe(prefix+"/connections/offer", func(self *mqtt.Client, message mqtt.Message) {
		var req ConnectionDescriptionMessage
		if err := json.Unmarshal(message.Payload(), &req); err != nil {
			log.Println(err.Error())
			return
		}

		peer, answer := service.webRTCService.NewConnection(req.Body.SessionDescription, core.PeerHooks{
			OnCandidate: func(p core.Peer, i webrtc.ICECandidateInit) {
				message := ConnectionCandidateMessage{
					Headers: Headers{
						SenderDeviceId:  prefix,
						SenderSessionId: p.Id,
						SessionId:       req.Headers.SenderSessionId,
					},
					Body: Body{
						IceCandidate: i,
					},
				}
				payload, _ := json.Marshal(message)
				self.Publish(req.Headers.SenderDeviceId+"/connections/candidate", string(payload))
			},
			OnClose: func(p core.Peer) {
				message := ConnectionCloseMessage{
					Headers: Headers{
						SenderDeviceId:  prefix,
						SenderSessionId: p.Id,
						SessionId:       req.Headers.SenderSessionId,
					},
				}
				payload, _ := json.Marshal(message)
				self.Publish(req.Headers.SenderDeviceId+"/connections/close", string(payload))
			},
		})

		response := ConnectionDescriptionMessage{
			Headers: Headers{
				SenderDeviceId:  prefix,
				SenderSessionId: peer.Id,
				SessionId:       req.Headers.SenderSessionId,
			},
			Body: Body{
				SessionDescription: answer,
			},
		}
		payload, _ := json.Marshal(response)
		client.Publish(req.Headers.SenderDeviceId+"/connections/answer", string(payload))
	})

	client.Subscribe(prefix+"/connections/close", func(self *mqtt.Client, message mqtt.Message) {
		var req ConnectionCloseMessage
		if err := json.Unmarshal(message.Payload(), &req); err != nil {
			log.Println(err.Error())
			return
		}

		peer := service.webRTCService.GetConnectionById(req.Headers.SessionId)
		if peer == nil {
			return
		}

		service.webRTCService.CloseConnection(peer)
	})

	client.Subscribe(prefix+"/connections/candidate", func(self *mqtt.Client, message mqtt.Message) {
		var req ConnectionCandidateMessage
		if err := json.Unmarshal(message.Payload(), &req); err != nil {
			log.Println(err.Error())
			return
		}

		peer := service.webRTCService.GetConnectionById(req.Headers.SessionId)
		if peer == nil {
			return
		}

		service.webRTCService.AddICECandidate(peer, req.Body.IceCandidate)
	})

	client.Subscribe(prefix+"/actions/trigger", func(self *mqtt.Client, message mqtt.Message) {
		var req struct {
			Pattern     string            `json:"pattern"`
			Environment map[string]string `json:"environment"`
		}
		if err := json.Unmarshal(message.Payload(), &req); err != nil {
			log.Println(err.Error())
			return
		}

		service.actionService.Execute(req.Pattern, req.Environment)
	})
}
