package core

import "github.com/dieklingel/core/internal/mqtt"

type MqttConnection struct {
	Id       uint64       `gorm:"primaryKey"`
	Client   *mqtt.Client `gorm:"-"`
	Url      string
	Username string
	Password string
}

type MqttService interface {
	Connections() []MqttConnection
	GetConnectionById(id int) *MqttConnection
	SaveConnection(connection *MqttConnection)
	RemoveConnection(connection *MqttConnection)
}
