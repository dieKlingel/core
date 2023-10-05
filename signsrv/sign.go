package signsrv

import (
	"errors"
	"log"

	"github.com/dieklingel/core/internal/core"
	"gorm.io/gorm"
)

type SignService struct {
	database              *gorm.DB
	onSignSavedHandlers   []func(core.Sign)
	onSignRemovedHandlers []func(core.Sign)
}

func NewService(db *gorm.DB) core.SignService {
	db.AutoMigrate(&core.Sign{})

	return &SignService{
		database:              db,
		onSignSavedHandlers:   make([]func(core.Sign), 0),
		onSignRemovedHandlers: make([]func(core.Sign), 0),
	}
}

func (service *SignService) Signs() []core.Sign {
	var signs []core.Sign
	if res := service.database.Find(&signs); res.Error != nil {
		log.Print(res.Error.Error())
	}

	return signs
}

func (service *SignService) GetSignById(id int) *core.Sign {
	var sign core.Sign
	err := service.database.First(&sign, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &sign
}

func (service *SignService) SaveSign(sign *core.Sign) {
	service.database.Save(sign)
	go func(sign core.Sign) {
		for _, handler := range service.onSignSavedHandlers {
			go handler(sign)
		}
	}(*sign)
}

func (service *SignService) RemoveSign(sign *core.Sign) {
	service.database.Delete(sign)
	go func(sign core.Sign) {
		for _, handler := range service.onSignRemovedHandlers {
			go handler(sign)
		}
	}(*sign)
}

func (service *SignService) OnSignSaved(handler func(sign core.Sign)) {
	service.onSignSavedHandlers = append(service.onSignSavedHandlers, handler)
}

func (service *SignService) OnSignRemoved(handler func(sign core.Sign)) {
	service.onSignRemovedHandlers = append(service.onSignRemovedHandlers, handler)
}
