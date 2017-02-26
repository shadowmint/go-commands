# go-commands

Actions is a generic framework for task based actions between modules.

The idea is that you create an `Commands` object, then register task handlers
on it, using `Register`.

Then other modules can come along an request that a particular type of action
be dispatched in a generic way.

# Usage

    npm install shadowmint/go-commands --save

You can use command sync:

    instance := commands.New()
    instance.Register(&FooCommandHandler{})

    cmd := newFooCommand()
    cmd.EventHandler().Listen(commands.CommandCompleted{}, func(c interface{}) {
        ...
    })

    err := instance.Wait(cmd)

Or async:

    promise, err := instance.Execute(cmd)
    promise.Then(func() { ... }, func(err error) { ... })

If you're using `go-bind` you can supply a registry object to use:

    type FooCommandHandler struct {
        Api ServiceApi
    }

	reg := threaded.New()
    reg.Register((*ServiceApi)(nil), func(r registry.Registry) (interface{}, error) {
        return &ServiceImpl{}, nil
    })

	instance := commands.New(reg)
	instance.Register(&FooCommandHandler{})

The command handler will be resolved using the registry instance.