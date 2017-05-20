package commands

import (
	"fmt"
	"reflect"
	"time"

	"ntoolkit/errors"
	"ntoolkit/futures"
	"ntoolkit/registry"
	"sync"
)

// Commands is a high level construct for dispatching commands and registering command handlers.
type Commands struct {
	// Timeout is the default command timeout for long running commands
	Timeout  *time.Duration
	handlers map[reflect.Type]CommandHandler
	registry registry.Registry
}

// New returns a new instance of the Commands type; if a registry is supplied
// it is used to bind command handlers before they get executed.
func New(registry ...registry.Registry) *Commands {
	rtn := &Commands{
		Timeout:  nil,
		handlers: make(map[reflect.Type]CommandHandler)}
	if len(registry) > 0 {
		rtn.registry = registry[0]
	}
	return rtn
}

// Register a new command handler
func (c *Commands) Register(handler CommandHandler) error {
	if handler == nil {
		return errors.Fail(ErrBadHandler{}, nil, "nil is not a valid command handler")
	}
	if c.registry != nil {
		if err := c.registry.Bind(handler); err != nil {
			return err
		}
	}
	c.handlers[handler.Handles()] = handler
	return nil
}

// Wait for a command to finish and return an error if it fails
func (c *Commands) Wait(command Command) error {
	var err error
	promise := c.Execute(command)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		promise.Then(func() {
			wg.Done()
		}, func(cmdErr error) {
			err = cmdErr
			wg.Done()
		})
	}()
	wg.Wait()
	return err
}

// Execute a command and return an error if it failed
func (c *Commands) Execute(command Command) *futures.Deferred {
	rtn := &futures.Deferred{}
	if command == nil {
		rtn.Reject(errors.Fail(ErrNoHandler{}, nil, "No command handler for nil"))
		return rtn
	}

	if handler, ok := c.handlers[reflect.TypeOf(command)]; ok {
		if setup, ok := command.(Setup); ok {
			setup.Setup()
		}

		// Timeout?
		timeout := c.Timeout
		if timeoutRef, ok := command.(Timeout); ok {
			timeout = timeoutRef.Timeout()
		}

		return c.executeTimed(command, handler, timeout)
	}

	rtn.Reject(errors.Fail(ErrNoHandler{}, nil, fmt.Sprintf("No command handler found for unknown type: %s", reflect.TypeOf(command))))
	return rtn
}

// executeTimed returns a promise for command execution that invokes the
// required interfaces and return a promise for being finished.
func (c *Commands) executeTimed(command Command, handler CommandHandler, timeout *time.Duration) *futures.Deferred {
	rtn := futures.Deferred{}
	lock := newLockedBool()
	promise := handler.Execute(command)
	promise.Then(func() {
		lock.Enter(func() error {
			if completed, ok := command.(Completed); ok {
				completed.Completed()
			}
			command.EventHandler().Trigger(CommandCompleted{command, nil})
			rtn.Resolve()
			return nil
		})
	}, func(err error) {
		lock.Enter(func() error {
			if failed, ok := command.(Failed); ok {
				failed.Failed(err)
			}
			command.EventHandler().Trigger(CommandCompleted{command, err})
			rtn.Reject(err)
			return nil
		})
		rtn.Reject(err)
	})
	if timeout != nil {
		go (func() {
			time.Sleep(*timeout)
			lock.Enter(func() error {
				err := errors.Fail(ErrBadHandler{}, nil, fmt.Sprintf("Timeout after %d ms running command", timeout))
				command.EventHandler().Trigger(CommandCompleted{command, err})
				rtn.Reject(err)
				return nil
			})
		})()
	}
	return &rtn
}
