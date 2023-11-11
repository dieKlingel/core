package api

type Action interface {
	Id() string
	Trigger() string
	Script() string
}

type ActionExecutionResult interface {
	Action() Action
	ExitCode() int
	Output() string
}
