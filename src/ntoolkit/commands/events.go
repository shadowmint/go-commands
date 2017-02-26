package commands

type CommandCompleted struct {
	Command Command
	Error error
}