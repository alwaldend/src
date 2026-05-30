package al

import (
	"errors"
)

// Panic on errors
func Check(errs ...error) {
	err := errors.Join(errs...)
	if err != nil {
		panic(err)
	}
}

func Must[T any](val T, err error) T {
	Check(err)
	return val
}
