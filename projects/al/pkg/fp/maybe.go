package fp

import "fmt"

type Maybe[T any] struct {
	value *T
}

var _ Monad[int] = (*Maybe[int])(nil)
var _ Functor[int] = (*Maybe[int])(nil)

func Just[T any](x T) Maybe[T] {
	return Maybe[T]{&x}
}

func Nothing[T any]() Maybe[T] {
	return Maybe[T]{nil}
}

func (self Maybe[T]) Get() (T, error) {
	if self.value == nil {
		var res T
		return res, fmt.Errorf("maybe is missing value")
	}
	return *self.value, nil
}

func (self Maybe[T]) Map(f func(T) T) Functor[T] {
	if self.value == nil {
		return Nothing[T]()
	}
	return Just(f(*self.value))
}

func (self Maybe[T]) FlatMap(f func(T) Monad[T]) Monad[T] {
	if self.value == nil {
		return Nothing[T]()
	}
	return f(*self.value)
}
