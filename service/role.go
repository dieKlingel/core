package service

import (
	"encoding/json"
	"log"

	"github.com/dieklingel/core/internal/role"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

type RoleService struct{}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func (service *RoleService) RuleSet() *role.RuleSet {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return role.NewRuleSet()
	}
	defer db.Close()

	if succes, _ := db.HasCollection("roles"); !succes {
		return role.NewRuleSet()
	}

	doc, err := db.FindFirst(query.NewQuery("roles"))
	if err != nil {
		log.Printf("an error occured while fetching actions: %s", err.Error())
		return role.NewRuleSet()
	}

	if doc == nil {
		return role.NewRuleSet()
	}

	policy := make(map[string]map[string]bool)
	if err = json.Unmarshal([]byte(doc.Get("policy").(string)), &policy); err != nil {
		return role.NewRuleSet()
	}

	return role.NewRuleSetFromMap(policy)
}

func (service *RoleService) SetRuleSet(rules *role.RuleSet) error {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return err
	}
	defer db.Close()

	if succes, _ := db.HasCollection("roles"); !succes {
		db.CreateCollection("roles")
	}

	policy, err := json.Marshal(rules.ToMap())
	if err != nil {
		return err
	}

	doc, err := db.FindFirst(query.NewQuery("roles"))
	if err != nil {
		return err
	}

	if doc == nil {
		doc = document.NewDocument()
	}
	doc.Set("policy", string(policy))

	return db.Save("roles", doc)
}
