package api

type User interface {
	Username() string
	Password() string
	Role() string
}
