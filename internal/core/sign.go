package core

type Sign interface {
	Id() string
	Script() string
}

type SignService interface {
	Signs() []Sign
	SaveSign(sign Sign)
	RemoveSign(sign Sign)
	OnSignSaved(handler func(sign Sign))
	OnSignRemoved(handler func(sign Sign))
}
