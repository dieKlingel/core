package core

type Sign struct {
	Id     int
	Script string
}

type SignService interface {
	Signs() []Sign
	GetSignById(id int) *Sign
	SaveSign(sign *Sign)
	RemoveSign(sign *Sign)
	OnSignSaved(handler func(sign Sign))
	OnSignRemoved(handler func(sign Sign))
}
