package core

type User struct {
	Username     string `gorm:"primaryKey"`
	PasswordHash string
}

type UserService interface {
	Users() []User
	GetUserByUsername(username string) *User
	SaveUser(user *User)
	RemoveUser(user *User)
	OnUserSaved(handler func(user User))
	OnUserRemoved(handler func(user User))
}
