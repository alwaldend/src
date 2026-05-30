package fp

import (
	"fmt"
	"os"
)

type Result[T any] interface {
	Get() (T, error)
}

func Get[T any](val Functor[T]) (res T, err error) {
	valTyped, ok := val.(Result[T])
	if !ok {
		return res, fmt.Errorf("could not convert %T to ResultContainer", val)
	}
	return valTyped.Get()
}

func GetOrExit[T any](val Functor[T], code int) T {
	res, err := Get(val)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(code)
	}
	return res
}

func GetOrPanic[T any](val Functor[T]) T {
	res, err := Get(val)
	if err != nil {
		panic(err)
	}
	return res
}
