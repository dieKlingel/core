package usersrv

import (
	"errors"
	"log"

	"github.com/dieklingel/core/internal/core"
	"gorm.io/gorm"
)

type UserService struct {
	database              *gorm.DB
	onUserSavedHandlers   []func(core.User)
	onUserRemovedHandlers []func(core.User)
}

func NewService(db *gorm.DB) core.UserService {
	db.AutoMigrate(&core.User{})

	return &UserService{
		database:              db,
		onUserSavedHandlers:   make([]func(core.User), 0),
		onUserRemovedHandlers: make([]func(core.User), 0),
	}
}

func (service *UserService) Users() []core.User {
	var users []core.User
	if res := service.database.Find(&users); res.Error != nil {
		log.Print(res.Error.Error())
	}

	return users
}

func (service *UserService) GetUserByUsername(username string) *core.User {
	var user core.User
	err := service.database.First(&user, "username = ?", username).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &user
}

func (service *UserService) SaveUser(user *core.User) {
	service.database.Save(user)
}

func (service *UserService) RemoveUser(user *core.User) {
	service.database.Delete(user)
}

func (service *UserService) OnUserSaved(handler func(core.User)) {

}

func (service *UserService) OnUserRemoved(handler func(core.User)) {

}
