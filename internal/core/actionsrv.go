package core

type ActionExecutionEnvironment string

const (
	ActionExecutionEnvironmentBash   ActionExecutionEnvironment = "bash"
	ActionExecutionEnvironmentPython ActionExecutionEnvironment = "python"
)

type ActionHandler func(env map[string]string)

type Action struct {
	Trigger     string
	Script      string
	Environment ActionExecutionEnvironment
}

type ActionService interface {
	Register(trigger string, handler ActionHandler)
	Execute(pattern string, environment map[string]string)
}
