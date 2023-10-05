package core

type User interface {
	Id() string
	Role() Role
}

type UserService interface {
	Users() []User
	SaveUser(user User)
	RemoveUser(user User)
	OnUserSaved(handler func(user User))
	OnUserRemoved(handler func(user User))
}
