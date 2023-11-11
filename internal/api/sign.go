package api

type Sign interface {
	Id() string
	Name() string
	Script() string
}
