package fp

import "fmt"

type Maybe[T any] struct {
	value *T
}

var _ Monad[Maybe[int], int] = (*Maybe[int])(nil)
var _ Functor[Maybe[int], int] = (*Maybe[int])(nil)
var _ Result[int] = (*Maybe[int])(nil)

func Just[T any](x T) Maybe[T] {
	return Maybe[T]{&x}
}

func Nothing[T any]() Maybe[T] {
	return Maybe[T]{nil}
}

func (self Maybe[T]) Get() (T, error) {
	if self.value == nil {
		return *new(T), fmt.Errorf("maybe is missing value")
	}
	return *self.value, nil
}

func (self Maybe[T]) Map(f func(T) T) Maybe[T] {
	if self.value == nil {
		return Nothing[T]()
	}
	return Just(f(*self.value))
}

func (self Maybe[T]) FlatMap(f func(T) Maybe[T]) Maybe[T] {
	if self.value == nil {
		return Nothing[T]()
	}
	return f(*self.value)
}
