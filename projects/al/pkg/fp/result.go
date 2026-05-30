package fp

import "fmt"

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
