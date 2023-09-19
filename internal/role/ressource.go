package role

type Ressource struct {
	name string
	role *Role
}

func (res *Ressource) Allow() {
	res.Set(true)
}

func (res *Ressource) Deny() {
	res.Set(false)
}

func (res *Ressource) Set(allow bool) {
	res.role.ruleset.setPermission(res.role.name, res.name, allow)
}

func (res *Ressource) Get() bool {
	return res.role.ruleset.permission(res.role.name, res.name)
}

func (res *Ressource) Unset() {
	res.role.ruleset.resetPermission(res.role.name, res.name)
}
