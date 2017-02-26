package commands

import (
	"fmt"
	"sync"

	"ntoolkit/errors"
)

type lockedBool struct {
	state bool
	lock  *sync.Mutex
}

func newLockedBool() *lockedBool {
	return &lockedBool{false, &sync.Mutex{}}
}

// Enter locks the internal mutex and executes the action given if its false.
// If task wasn't executed or failed, an error is returned.
func (lock *lockedBool) Enter(task func() error) error {
	lock.lock.Lock()
	defer (func() {
		lock.lock.Unlock()
	})()
	if lock.state {
		return errors.Fail(ErrBadHandler{}, nil, "Unable to lock command for update")
	}
	return lock.execute(task)
}

// execute runs a task and return an error if it fails
func (lock *lockedBool) execute(task func() error) (err error) {
	defer (func() {
		r := recover()
		if r != nil {
			if evalue, ok := r.(error); ok {
				err = evalue
			} else {
				err = errors.Fail(ErrBadHandler{}, nil, fmt.Sprintf("Unknown failure to execute task: %s", r))
			}
		}
	})()
	return task()
}
