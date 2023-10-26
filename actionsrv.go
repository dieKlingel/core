package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"github.com/dieklingel/core/internal/core"
)

type ActionService struct {
	storageService core.StorageService
	handlers       map[string][]ActionHandler
}

type ActionHandler func(env map[string]string)

func NewActionService(storageService core.StorageService) *ActionService {
	return &ActionService{
		storageService: storageService,
	}
}

func (actionService *ActionService) Register(trigger string, handler ActionHandler) {
	if _, exists := actionService.handlers[trigger]; !exists {
		actionService.handlers[trigger] = make([]ActionHandler, 0)
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

	actions := actionService.storageService.Read().Actions
	for _, action := range actions {
		var command *exec.Cmd

		switch action.Environment {
		case core.ActionExecutionEnvironmentBash:
			command = exec.Command("bash", "-c", action.Script)
		case core.ActionExecutionEnvironmentPython:
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
		log.Println(output)
	}
}
