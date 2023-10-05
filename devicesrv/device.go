package devicesrv

import (
	"errors"
	"log"

	"github.com/dieklingel/core/internal/core"
	"gorm.io/gorm"
)

type DeviceService struct {
	database                *gorm.DB
	onDeviceSavedHandlers   []func(core.Device)
	onDeviceRemovedHandlers []func(core.Device)
}

func NewService(db *gorm.DB) core.DeviceService {
	db.AutoMigrate(&core.Device{})

	return &DeviceService{
		database:                db,
		onDeviceSavedHandlers:   make([]func(core.Device), 0),
		onDeviceRemovedHandlers: make([]func(core.Device), 0),
	}
}

func (service *DeviceService) Devices() []core.Device {
	var devices []core.Device
	if res := service.database.Find(&devices); res.Error != nil {
		log.Print(res.Error.Error())
	}

	return devices
}

func (service *DeviceService) GetDeviceById(id int) *core.Device {
	var device core.Device
	err := service.database.First(&device, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &device
}

func (service *DeviceService) SaveDevice(device *core.Device) {
	service.database.Save(device)
}

func (service *DeviceService) RemoveDevice(device *core.Device) {
	service.database.Delete(device)
}

func (service *DeviceService) OnDeviceSaved(handler func(device core.Device)) {
	service.onDeviceSavedHandlers = append(service.onDeviceSavedHandlers, handler)
}

func (service *DeviceService) OnDeviceRemoved(handler func(device core.Device)) {
	service.onDeviceRemovedHandlers = append(service.onDeviceRemovedHandlers, handler)
}
