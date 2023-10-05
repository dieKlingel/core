package core

type Role interface {
	Id() string
}

type RoleService interface {
	Roles() []Role
	SaveRole(role Role)
	RemoveRole(role Role)
	OnRoleSaved(handler func(role Role))
	OnRoleRemoved(handler func(role Role))
}
