package endpoint

import "github.com/dieklingel/core/internal/role"

type RoleService interface {
	RuleSet() *role.RuleSet
}

type RoleEndpoint struct {
	service RoleService
}

func NewRoleEndpoint(service RoleService) *RoleEndpoint {
	return &RoleEndpoint{
		service: service,
	}
}

func (endpoint *RoleEndpoint) RuleSet() *role.RuleSet {
	return endpoint.service.RuleSet()
}
