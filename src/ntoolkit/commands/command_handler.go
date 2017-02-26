package commands

import (
	"reflect"
	"ntoolkit/futures"
)

type CommandHandler interface {
	// Handles returns the type supported by this command handler
	Handles() reflect.Type

	// Execute executes the command given and returns an error on failure
	Execute(command interface{}) *futures.Deferred
}
