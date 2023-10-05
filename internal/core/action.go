package core

import (
	"fmt"
	"os/exec"
)

type ActionExecutionEnvironment string

const (
	ActionExecutionEnvironmentBash   ActionExecutionEnvironment = "bash"
	ActionExecutionEnvironmentPython ActionExecutionEnvironment = "python"
)

type ActionExecutionResult struct {
	ErrorMessage string
	Output       string
	ActionId     uint64
}

type Action struct {
	Id          uint64 `gorm:"primaryKey"`
	Trigger     string
	Script      string
	Environment ActionExecutionEnvironment
}

func (action *Action) Execute(env map[string]string) ActionExecutionResult {
	result := ActionExecutionResult{
		ActionId: action.Id,
	}
	var command *exec.Cmd

	switch action.Environment {
	case ActionExecutionEnvironmentBash:
		command = exec.Command("bash", "-c", action.Script)
	case ActionExecutionEnvironmentPython:
		command = exec.Command("python3", "-c", action.Script)
	default:
		result.ErrorMessage = fmt.Sprintf("the action %d has not defined a supported execution environment", action.Id)
		return result
	}

	for key, value := range env {
		command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, value))
	}

	output, err := command.CombinedOutput()
	if err != nil {
		result.ErrorMessage = err.Error()
	}
	result.Output = string(output)
	return result
}

type ActionService interface {
	Actions() []Action
	GetActionById(id int) *Action
	SaveAction(action Action) Action
	RemoveAction(action Action) Action
	OnActionSaved(handler func(action Action))
	OnActionRemoved(handler func(action Action))
	Execute(pattern string, env map[string]string) []ActionExecutionResult
}
