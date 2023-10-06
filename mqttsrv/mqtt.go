package mqttsrv

import (
	"github.com/dieklingel/core/internal/core"
	"github.com/dieklingel/core/internal/mqtt"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type MqttService struct {
	database   *gorm.DB
	connetions []core.MqttConnection
}

func NewService(db *gorm.DB) core.MqttService {
	db.AutoMigrate(&core.MqttConnection{})

	var connections []core.MqttConnection
	db.Find(&connections)

	service := &MqttService{
		database:   db,
		connetions: connections,
	}

	for index := range service.connetions {
		service.connect(&service.connetions[index])
	}

	return service
}

func (service *MqttService) Connections() []core.MqttConnection {
	return service.connetions
}

func (service *MqttService) GetConnectionById(id int) *core.MqttConnection {
	for index := range service.connetions {
		if service.connetions[index].Id == uint64(id) {
			return &service.connetions[index]
		}
	}
	return nil
}

func (service *MqttService) SaveConnection(connection *core.MqttConnection) {
	service.database.Save(connection)
	if existingConn := service.GetConnectionById(int(connection.Id)); existingConn != nil {
		service.connect(existingConn)
		connection.Client = existingConn.Client
	} else {
		service.connect(connection)
		service.connetions = append(service.connetions, *connection)
	}
}

func (service *MqttService) RemoveConnection(connection *core.MqttConnection) {
	service.database.Delete(connection)
	if connection.Client != nil {
		connection.Client.Disconnect()
	}

	index := slices.Index[[]core.MqttConnection](service.connetions, *connection)
	if index == -1 {
		return
	}

	service.connetions = slices.Delete(service.connetions, index, index+1) //append(service.connetions[:index], service.connetions[index+1:]...)
}

func (service *MqttService) connect(connection *core.MqttConnection) {
	if connection.Client != nil {
		connection.Client.Disconnect()
	}
	connection.Client = mqtt.NewClient()
	connection.Client.SetAutoReconnect(true)
	connection.Client.SetBroker(connection.Url)
	connection.Client.SetUsername(connection.Username)
	connection.Client.SetPassword(connection.Password)
	connection.Client.Connect()
}
