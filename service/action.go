package service

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"github.com/dieklingel/core/internal/api"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

type action struct {
	id      string
	trigger string
	script  string
}

func (a *action) Id() string {
	return a.id
}

func (a *action) Trigger() string {
	return a.trigger
}

func (a *action) Script() string {
	return a.script
}

type actionExecutionResult struct {
	action   api.Action
	exitCode int
	output   string
}

func (res *actionExecutionResult) Action() api.Action {
	return res.action
}

func (res *actionExecutionResult) ExitCode() int {
	return res.exitCode
}

func (res *actionExecutionResult) Output() string {
	return res.output
}

type ActionService struct{}

func NewActionService() *ActionService {
	return &ActionService{}
}

func (service *ActionService) List() []api.Action {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return make([]api.Action, 0)
	}
	defer db.Close()

	if succes, _ := db.HasCollection("actions"); !succes {
		return make([]api.Action, 0)
	}

	docs, err := db.FindAll(query.NewQuery("actions"))
	if err != nil {
		log.Printf("an error occured while fetching actions: %s", err.Error())
		return make([]api.Action, 0)
	}

	actions := make([]api.Action, len(docs))
	for index, doc := range docs {
		id := doc.ObjectId()
		trigger := doc.Get("trigger").(string)
		script := doc.Get("script").(string)

		actions[index] = &action{
			id:      id,
			trigger: trigger,
			script:  script,
		}
	}

	return actions
}

func (service *ActionService) Filter(pattern string) []api.Action {
	actions := service.List()
	filterd := make([]api.Action, 0)

	for _, action := range actions {
		if matched, _ := regexp.MatchString(pattern, action.Trigger()); matched {
			filterd = append(filterd, action)
		}
	}

	return filterd
}

func (service *ActionService) Remove(id string) error {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return err
	}
	defer db.Close()

	if succes, _ := db.HasCollection("actions"); !succes {
		return nil
	}

	return db.DeleteById("actions", id)
}

func (service *ActionService) Add(trigger string, script string) (string, error) {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return "", err
	}
	defer db.Close()

	if succes, _ := db.HasCollection("actions"); !succes {
		db.CreateCollection("actions")
	}

	doc := document.NewDocument()
	doc.Set("trigger", trigger)
	doc.Set("script", script)

	if err := db.Insert("actions", doc); err != nil {
		log.Printf("an error occured while inserting into the actions database: %s", err.Error())
		return "", err
	}

	return doc.ObjectId(), nil
}

func (service *ActionService) Execute(action api.Action, environment map[string]string) api.ActionExecutionResult {
	command := exec.Command("bash", "-c", action.Script())

	for key, value := range environment {
		command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, value))
	}

	output, _ := command.CombinedOutput()

	return &actionExecutionResult{
		action:   action,
		output:   string(output),
		exitCode: -1, // TODO: read exit code
	}
}
