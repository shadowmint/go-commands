package commands_test

import (
	"reflect"
	"ntoolkit/events"
	"ntoolkit/futures"
)

type FooCommand struct {
	handler *events.EventHandler
}

func newFooCommand() *FooCommand {
	return &FooCommand{handler: events.New()}
}

func (h *FooCommand) EventHandler() *events.EventHandler {
	return h.handler
}

type BarCommand struct {
	handler *events.EventHandler
	Success bool
	Pending bool
}

func newBarCommand() *BarCommand {
	return &BarCommand{handler: events.New()}
}

func (h *BarCommand) EventHandler() *events.EventHandler {
	return h.handler
}

func (h *BarCommand) Setup() {
	h.Pending = true
}

func (h *BarCommand) Completed() {
	h.Success = true
	h.Pending = false
}

func (h *BarCommand) Failed() {
	h.Success = false
	h.Pending = false
}

type FooCommandHandler struct {
	Api ServiceApi
	WithApi func(value bool)
}

func (h *FooCommandHandler) Handles() reflect.Type {
	return reflect.TypeOf(&FooCommand{})
}

func (h *FooCommandHandler) Execute(command interface{}) *futures.Deferred {
	rtn := &futures.Deferred{}
	if h.Api != nil {
		if h.WithApi != nil {
			h.WithApi(h.Api.Foo())
		}
	}
	rtn.Resolve()
	return rtn
}


type BarCommandHandler struct {
}

func (h *BarCommandHandler) Handles() reflect.Type {
	return reflect.TypeOf(&BarCommand{})
}

func (h *BarCommandHandler) Execute(command interface{}) *futures.Deferred {
	rtn := &futures.Deferred{}
	rtn.Resolve()
	return rtn
}

type ServiceImpl struct {
}

func (s *ServiceImpl) Foo() bool {
	return true
}

type ServiceApi interface {
	Foo() bool
}
