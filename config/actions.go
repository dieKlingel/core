package config

type ActionExecutionEnvironment string

const (
	ActionExecutionEnvironmentBash   ActionExecutionEnvironment = "bash"
	ActionExecutionEnvironmentPython ActionExecutionEnvironment = "python"
)

type Action struct {
	Trigger     string                     `yaml:"trigger"`
	Environment ActionExecutionEnvironment `yaml:"environment"`
	Script      string                     `yaml:"script"`
}
