package endpoint

import (
	"github.com/dieklingel/core/internal/api"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(username string, password string, role string) (api.User, error)
	GetByUsername(username string) api.User
}

type UserEndpoint struct {
	service UserService
	role    RoleService
}

func NewUserEndpoint(service UserService, role RoleService) *UserEndpoint {
	return &UserEndpoint{
		service: service,
		role:    role,
	}
}

func (endpoint *UserEndpoint) Create(username string, password string, role string) (api.User, error) {
	return endpoint.service.Create(username, password, role)
}

func (endpoint *UserEndpoint) GetByUsername(username string) api.User {
	return endpoint.service.GetByUsername(username)
}

func (endpoint *UserEndpoint) Authenticate(username string, password string) (bool, api.User) {
	user := endpoint.GetByUsername(username)
	if user == nil {
		return false, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password()), []byte(password)); err != nil {
		return false, nil
	}

	return true, user
}

func (endpoint *UserEndpoint) Authorize(username string, password string, ressource string) (bool, api.User) {
	_, user := endpoint.Authenticate(username, password)
	if user == nil {
		return false, nil
	}

	ruleset := endpoint.role.RuleSet()
	return ruleset.Role(user.Role()).Ressource(ressource).Get(), user
}
