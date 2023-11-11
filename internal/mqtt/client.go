package mqtt

import (
	"time"

	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type Message mq.Message

type Client struct {
	broker        string
	clientId      string
	username      string
	password      string
	autoReconnect bool
	keepAlive     time.Duration
	client        mq.Client
	error         error
}

func NewClient() *Client {
	client := &Client{
		broker:        "",
		clientId:      uuid.New().String(),
		username:      "",
		password:      "",
		autoReconnect: false,
		keepAlive:     0 * time.Second,
		client:        nil,
		error:         nil,
	}

	return client
}

func (client *Client) Error() error {
	return client.error
}

func (client *Client) ErrorMessage() string {
	if client.error == nil {
		return ""
	}
	return client.Error().Error()
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
		client.client.Disconnect(0)
	}

	options := mq.NewClientOptions()
	options.AddBroker(client.broker)
	options.SetClientID(client.clientId)
	options.SetUsername(client.username)
	options.SetPassword(client.password)
	options.SetAutoReconnect(client.autoReconnect)
	options.SetKeepAlive(client.keepAlive)

	c := mq.NewClient(options)
	client.client = c

	token := client.client.Connect()

	if token.Wait(); token.Error() != nil {
		client.error = token.Error()
	} else {
		client.error = nil
	}

	return token
}

func (client *Client) Disconnect() {
	if client.client == nil {
		return
	}

	client.client.Disconnect(0)
	client.client = nil
}

func (client *Client) IsConnected() bool {
	return client.client.IsConnected()
}

func (client *Client) Subscribe(topic string, handler func(self *Client, message Message)) {
	client.client.Subscribe(topic, 2, func(c mq.Client, m mq.Message) {
		handler(client, m)
	})
}

func (client *Client) Publish(topic string, message string) {
	client.client.Publish(topic, 2, false, message)
}
