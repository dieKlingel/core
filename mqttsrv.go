package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dieklingel/core/internal/core"
	"github.com/dieklingel/core/internal/mqtt"
	"github.com/pion/webrtc/v3"
)

type MqttService struct {
	storageService core.StorageService
	actionService  core.ActionService
	webRTCService  *WebRTCService
	appService     core.AppService

	client *mqtt.Client
}

func NewMqttService(storageService core.StorageService, actionsrv core.ActionService, webrtcsrc *WebRTCService, appService core.AppService) *MqttService {
	return &MqttService{
		storageService: storageService,
		actionService:  actionsrv,
		webRTCService:  webrtcsrc,
		appService:     appService,
	}
}

func (mqttService *MqttService) Run() {
	if mqttService.client != nil {
		mqttService.client.Disconnect()
	}

	config := mqttService.storageService.Read()
	url := config.Mqtt.Server
	username := config.Mqtt.Username
	password := config.Mqtt.Password

	mqttService.client = mqtt.NewClient()
	mqttService.client.SetAutoReconnect(true)
	mqttService.client.SetBroker(url)
	mqttService.client.SetUsername(username)
	mqttService.client.SetPassword(password)
	go func() {
		mqttService.client.Connect()

		for !mqttService.client.IsConnected() {
			fmt.Printf("could not connect to %s; retry in 10 src", url)
			time.Sleep(10 * time.Second)
			mqttService.client.Connect()
		}

		mqttService.buildListeners(mqttService.client, "")
	}()
}

func (service *MqttService) buildListeners(client *mqtt.Client, prefix string) {
	type Headers struct {
		SenderDeviceId  string `json:"senderDeviceId"`
		SenderSessionId string `json:"senderSessionId"`
		SessionId       string `json:"sessionId"`
	}

	type Body struct {
		SessionDescription webrtc.SessionDescription `json:"sessionDescription,omitempty"`
		IceCandidate       webrtc.ICECandidateInit   `json:"iceCandidate,omitempty"`
	}

	type ConnectionDescriptionMessage struct {
		Headers Headers `json:"headers"`
		Body    Body    `json:"body,omitempty"`
	}

	type ConnectionCandidateMessage struct {
		Headers Headers `json:"headers"`
		Body    Body    `json:"body"`
	}

	type ConnectionCloseMessage struct {
		Headers Headers `json:"headers"`
	}

	client.Subscribe(prefix+"/connections/create", func(self *mqtt.Client, message mqtt.Message) {
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

	client.Subscribe(prefix+"/apps/add", func(self *mqtt.Client, message mqtt.Message) {
		type Request struct {
			Headers struct {
				SenderDeviceId string `json:"senderDeviceId"`
			} `json:"headers"`
			Body struct {
				Token string `json:"token"`
			} `json:"body"`
		}

		var req Request
		if err := json.Unmarshal(message.Payload(), &req); err != nil {
			fmt.Printf("could not parse the AppsAddMessage; error: %v", err)
			return
		}

		service.appService.Register(req.Headers.SenderDeviceId, req.Body.Token)
	})

	client.Subscribe(prefix+"/apps/remove", func(self *mqtt.Client, message mqtt.Message) {
		type Request struct {
			Headers struct {
				SenderDeviceId string `json:"senderDeviceId"`
			} `json:"headers"`
		}

		var req Request
		if err := json.Unmarshal(message.Payload(), &req); err != nil {
			fmt.Printf("could not parse the AppsRemoveMessage; error: %v", err)
			return
		}

		service.appService.Unregister(req.Headers.SenderDeviceId)
	})
}
