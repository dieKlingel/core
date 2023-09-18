package endpoint

type SystemService interface {
	Version() string
}

type SystemEndpoint struct {
	service SystemService
}

func NewSystemEndpoint(service SystemService) *SystemEndpoint {
	return &SystemEndpoint{
		service: service,
	}
}

func (system *SystemEndpoint) Version() string {
	return system.service.Version()
}
