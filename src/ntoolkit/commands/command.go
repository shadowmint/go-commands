package commands

import (
	"time"
	"ntoolkit/events"
)

// Command is the base type an object must implement to be executable.
type Command interface {
	EventHandler() *events.EventHandler
}

// Setup is a pre-execute step for a command.
type Setup interface {
	Setup() error
}

// Completed is a post-execute step for commands that are successful.
type Completed interface {
	Completed()
}

// Failed is a post-execute step for commands that are not successful.
type Failed interface {
	Failed(err error)
}

// Timeout returns the maximum time in ms a task is permitted to run for before it is auto-failed.
// To allow tasks to wait forever, return nil; if not implemented the Commands waits for the default
// timeout value.
type Timeout interface {
	Timeout() *time.Duration
}
