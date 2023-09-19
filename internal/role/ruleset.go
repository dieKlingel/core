package role

type RuleSet struct {
	rules map[string]map[string]bool
}

func NewRuleSet() *RuleSet {
	return &RuleSet{
		rules: make(map[string]map[string]bool),
	}
}

func NewRuleSetFromMap(rules map[string]map[string]bool) *RuleSet {
	return &RuleSet{
		rules: rules,
	}
}

func (set *RuleSet) ToMap() map[string]map[string]bool {
	return set.rules
}

func (set *RuleSet) Role(role string) *Role {
	return &Role{
		name:    role,
		ruleset: set,
	}
}

func (set *RuleSet) permission(role string, subject string) bool {
	if sub, exists := set.rules[role]; exists {
		if allowed, exists := sub[subject]; exists {
			return allowed
		} else if wildcard, exists := sub["*"]; exists {
			return wildcard
		}
	}

	return false
}

func (set *RuleSet) setPermission(role string, subject string, allowed bool) {
	if _, exists := set.rules[role]; !exists {
		set.rules[role] = make(map[string]bool)
	}

	set.rules[role][subject] = allowed
}

func (set *RuleSet) resetPermission(role string, subject string) {
	if sub, exists := set.rules[role]; exists {
		delete(sub, subject)
	}
}
