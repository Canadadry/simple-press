package trycatch

import (
	"fmt"
	"runtime/debug"
)

var ThrowErr = fmt.Errorf("panic")

func Catch(fn func() error) error {
	var err error
	try(fn, &err)
	return err
}

func try(fn func() error, outerr *error) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%w : %v from %s", ThrowErr, r, string(debug.Stack()))
			*outerr = err
		}
	}()

	err := fn()
	if err != nil {
		*outerr = err
	}
}
