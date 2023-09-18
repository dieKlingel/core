package endpoint

import (
	"github.com/dieklingel/core/internal/api"
)

type ActionService interface {
	List() []api.Action
	Filter(pattern string) []api.Action
	Remove(id string) error
	Add(trigger string, script string) (string, error)
	Execute(action api.Action, environment map[string]string) api.ActionExecutionResult
}

type ActionEndpoint struct {
	service ActionService
}

// NewActionEndpoint returns a newly created Endpoint which uses the provided
// ActionService for all interactions with Actions.
func NewActionEndpoint(service ActionService) *ActionEndpoint {
	return &ActionEndpoint{
		service: service,
	}
}

// List returns a slice of all available Actions from the ActionService
func (endpoint *ActionEndpoint) List() []api.Action {
	return endpoint.service.List()
}

// Execute passes reads all Actions from the ActionService which match the
// provided Regex Pattern and Executes them in the order the ActionService
// delivers them. After all of the matching Actions are executed, a slice
// containing the result of all of the Actions will be returned.
// The order of the results in the slice matches the order of execution of the Actions
func (endpoint *ActionEndpoint) Execute(pattern string, environment map[string]string) []api.ActionExecutionResult {
	actions := endpoint.service.Filter(pattern)
	result := make([]api.ActionExecutionResult, len(actions))

	for index, action := range actions {
		result[index] = endpoint.service.Execute(action, environment)
	}

	return result
}

func (endpoint *ActionEndpoint) Add(trigger string, script string) string {
	id, err := endpoint.service.Add(trigger, script)
	if err != nil {
		return ""
	}

	return id
}
