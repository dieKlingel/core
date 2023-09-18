package endpoint

import "github.com/dieklingel/core/internal/api"

type SignService interface {
	List() []api.Sign
	Add(name string, script string) (string, error)
	Remove(id string) error
	GetById(id string) api.Sign
}

type SignEndpoint struct {
	service SignService
}

func NewSignEndpoint(service SignService) *SignEndpoint {
	return &SignEndpoint{
		service: service,
	}
}

func (endpoint *SignEndpoint) List() []api.Sign {
	return endpoint.service.List()
}

func (endpoint *SignEndpoint) Add(name string, script string) api.Sign {
	id, err := endpoint.service.Add(name, script)
	if err != nil {
		return nil
	}

	return endpoint.service.GetById(id)
}

func (endpoint *SignEndpoint) Delete(sign api.Sign) {
	endpoint.service.Remove(sign.Id())
}

func (endpoint *SignEndpoint) GetById(id string) api.Sign {
	return endpoint.service.GetById(id)
}
