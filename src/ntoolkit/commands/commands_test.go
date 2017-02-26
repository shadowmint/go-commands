package commands_test

import (
	"testing"

	"ntoolkit/assert"
	"ntoolkit/commands"
	"ntoolkit/registry/simple"
	"ntoolkit/registry/threaded"
	"ntoolkit/registry"
)

func TestNew(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		T.Assert(commands.New() != nil)

		r := simple.New()
		T.Assert(commands.New(r) != nil)
	})
}

func TestBoundCommandHandlerInvoked(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		instance := commands.New()
		instance.Register(&FooCommandHandler{})

		completed := false
		cmd := newFooCommand()
		cmd.EventHandler().Listen(commands.CommandCompleted{}, func(c interface{}) {
			completed = true
		})

		err := instance.Wait(cmd)
		T.Assert(err == nil)

		T.Assert(completed)

		err = instance.Wait(newBarCommand())
		T.Assert(err != nil)
	})
}

func TestAsyncCommand(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		instance := commands.New()
		instance.Register(&FooCommandHandler{})
		instance.Register(&BarCommandHandler{})

		cmd := newBarCommand()

		promise, err := instance.Execute(cmd)
		T.Assert(err == nil)

		promise.Then(func() {
			T.Assert(cmd.Success)
		}, func(err error) {
			T.Unreachable()
		})
	})
}

func TestRegisteredCommand(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		reg := threaded.New()
		reg.Register((*ServiceApi)(nil), func(r registry.Registry) (interface{}, error) {
			return &ServiceImpl{}, nil
		})

		resolved := false
		value := false
		instance := commands.New(reg)

		handler := &FooCommandHandler{WithApi: func(incoming bool) {
			resolved = true
			value = incoming
		}}

		instance.Register(handler)
		T.Assert(handler.Api != nil)

		cmd := newFooCommand()

		err := instance.Wait(cmd)
		T.Assert(err == nil)

		T.Assert(resolved)
		T.Assert(value)
	})
}