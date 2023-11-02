package core

type AppService interface {
	Register(id string, token string)
	Unregister(id string)
	IsRegisterd(id string) bool
}
