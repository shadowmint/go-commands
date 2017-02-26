package commands

// ErrNoHandler is raised when attempting to run a command with no bound handler.
type ErrNoHandler struct{}

// ErrBadCommandHandler is raised for invalid command handlers
type ErrBadHandler struct{}