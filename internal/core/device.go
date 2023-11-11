package core

type Device struct {
	Id uint64 `gorm:"primaryKey"`
}

type DeviceService interface {
	Devices() []Device
	GetDeviceById(id int) *Device
	SaveDevice(device *Device)
	RemoveDevice(device *Device)
	OnDeviceSaved(cb func(device Device))
	OnDeviceRemoved(cb func(device Device))
}
