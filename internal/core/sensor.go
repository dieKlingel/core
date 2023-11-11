package core

type SensorService interface {
	OnMovement(handler func())
}
