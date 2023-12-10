package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"github.com/dieklingel/core/config"
	"github.com/dieklingel/core/internal/core"
)

type ActionService struct {
	config   *config.Environment
	handlers map[string][]core.ActionHandler
}

func NewActionService(config *config.Environment) *ActionService {
	return &ActionService{
		config: config,
	}
}

func (actionService *ActionService) Register(trigger string, handler core.ActionHandler) {
	if _, exists := actionService.handlers[trigger]; !exists {
		actionService.handlers[trigger] = make([]core.ActionHandler, 0)
	}

	actionService.handlers[trigger] = append(actionService.handlers[trigger], handler)
}

func (actionService *ActionService) Execute(pattern string, env map[string]string) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for trigger, handlers := range actionService.handlers {
		if regex.Match([]byte(trigger)) {
			for _, handler := range handlers {
				go handler(env)
			}
		}
	}

	actions := actionService.config.Actions
	for _, action := range actions {
		if !regex.Match([]byte(action.Trigger)) {
			continue
		}

		var command *exec.Cmd

		switch action.Environment {
		case config.ActionExecutionEnvironmentBash:
			command = exec.Command("bash", "-c", action.Script)
		case config.ActionExecutionEnvironmentPython:
			command = exec.Command("python3", "-c", action.Script)
		default:
			log.Printf("the execution environment %s is not supported", action.Environment)
			continue
		}

		for key, value := range env {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, value))
		}

		output, err := command.CombinedOutput()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println(string(output))
	}
}
