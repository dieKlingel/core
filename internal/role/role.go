package role

type Role struct {
	name    string
	ruleset *RuleSet
}

func (role *Role) Ressource(ressource string) *Ressource {
	return &Ressource{
		name: ressource,
		role: role,
	}
}
