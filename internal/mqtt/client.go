package mqtt

import (
	"time"

	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type Client struct {
	broker        string
	clientId      string
	username      string
	password      string
	autoReconnect bool
	keepAlive     time.Duration
	client        *mq.Client
}

func NewClient() Client {
	client := Client{
		broker:        "",
		clientId:      uuid.New().String(),
		username:      "",
		password:      "",
		autoReconnect: false,
		keepAlive:     0 * time.Second,
		client:        nil,
	}

	return client
}

func (client *Client) SetBroker(server string) {
	client.broker = server
}

func (client *Client) SetClientId(id string) {
	client.clientId = id
}

func (client *Client) SetUsername(username string) {
	client.username = username
}

func (client *Client) SetPassword(password string) {
	client.password = password
}

func (client *Client) SetAutoReconnect(reconnect bool) {
	client.autoReconnect = reconnect
}

func (client *Client) SetKeepAlive(keepAlive time.Duration) {
	client.keepAlive = keepAlive
}

func (client *Client) Connect() Token {
	if client.client != nil {
		(*client.client).Disconnect(0)
	}

	options := mq.NewClientOptions()
	options.AddBroker(client.broker)
	options.SetClientID(client.clientId)
	options.SetUsername(client.username)
	options.SetPassword(client.password)
	options.SetAutoReconnect(client.autoReconnect)
	options.SetKeepAlive(client.keepAlive)

	c := mq.NewClient(options)
	client.client = &c

	return (*client.client).Connect()
}

func (client *Client) Disconnect() {
	if client.client == nil {
		return
	}

	(*client.client).Disconnect(0)
	client.client = nil

}
