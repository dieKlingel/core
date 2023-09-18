package transport

import "github.com/dieklingel/core/internal/api"

type SystemEndpoint interface {
	Version() string
}

type ActionEndpoint interface {
	List() []api.Action
	Execute(pattern string, environment map[string]string) []api.ActionExecutionResult
	GetById(id string) api.Action
	Add(trigger string, script string) api.Action
	Delete(api.Action)
}

type SignEndpoint interface {
	List() []api.Sign
	Add(name string, script string) api.Sign
	Delete(api.Sign)
	GetById(id string) api.Sign
}
